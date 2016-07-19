package session

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	// "github.com/gebv/echo-session"
	session "github.com/echo-contrib/sessions"
	"github.com/gorilla/sessions"
	faderstore "github.com/inpime/fader/store"
	"github.com/labstack/echo"
	"github.com/yosssi/boltstore/store"
	"net/http"
	"time"
)

var db *bolt.DB
var DefaultSessionName = "fssession"
var SessionNameContextKey = "_session"
var DefaultSessionUser *faderstore.File

// var LoggerKey = "addons_session"

type Config struct {
	DB *bolt.DB

	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool

	SessionName string
	SecretKey   string
	BucketName  string
}

func (s Config) TransformToSessionConfig() store.Config {
	cfg := store.Config{}
	cfg.SessionOptions.Domain = s.Domain
	cfg.SessionOptions.HttpOnly = s.HttpOnly
	cfg.SessionOptions.MaxAge = s.MaxAge
	cfg.SessionOptions.Path = s.Path
	cfg.SessionOptions.Secure = s.Secure

	cfg.DBOptions.BucketName = []byte(s.BucketName)

	return cfg
}

// InitSession init session storage
// NOTICE: only supported boltdb
// NOTICE: by default the guest is users@guestuser (check the bucket UsersBucketName)
func InitSession() (err error) {
	// utils.EnsureDir(filepath.Dir(Cfg.Session.Store.BoltDBFilePath))

	db, err = bolt.Open("session.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return err
	}
	// Init Default guest user session
	file, err := faderstore.LoadOrNewFile("users", "guest")
	if err != nil {
		return err
	}

	DefaultSessionUser = file

	return
}

type Store struct {
	sessions.Store
}

func (s *Store) Options(opt session.Options) {}

func SessionStoreMiddleware(name string, config Config) echo.MiddlewareFunc {
	_store, err := store.New(config.DB, config.TransformToSessionConfig(), []byte(config.SecretKey))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"_api": addonName,
		}).
			WithError(err).
			Panic("init store session")
	}

	logrus.Error("init store session")

	return session.Sessions(name, &Store{_store})
}

func InitializerUserSessionMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			uri := ctx.Request().URI()

			internalSession := session.Default(ctx)

			if internalSession == nil {
				// TODO: clear session or panic?

				logrus.WithFields(logrus.Fields{
					"_api": addonName,
					"url":  uri,
				}).Fatal("current session is null")
				return ctx.NoContent(http.StatusInternalServerError)
			}
			//
			// Session current request
			//

			_session := NewSession(internalSession.(session.Session))
			ctx.Set(SessionNameContextKey, _session)

			return h(ctx)
		}
	}
}
