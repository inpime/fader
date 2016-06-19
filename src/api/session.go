package api

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/gebv/echo-session"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/yosssi/boltstore/store"
	"net/http"
	"path/filepath"
	faderstore "store"
	"time"
	"utils"
)

var sessionDb *bolt.DB

var DefaultSessionName string = "fsession"

var ctxKeyNameSession string = "session"

type SessionConfig struct {
	DB *bolt.DB

	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool

	SecretKey  string
	BucketName string
}

func (s SessionConfig) TransformToSessionConfig() store.Config {
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
	if Cfg.Session.Store.Provider != "boltdb" {
		return fmt.Errorf("not supported session store %s", Cfg.Session.Store)
	}

	utils.EnsureDir(filepath.Dir(Cfg.Session.Store.BoltDBFilePath))

	db, err := bolt.Open(Cfg.Session.Store.BoltDBFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return err
	}

	sessionDb = db

	// Init Default guest user session
	file, err := faderstore.LoadOrNewFile(UsersBucketName, GuestUserFileName)
	if err != nil {
		return err
	}

	DefaultGuestSession = file

	return
}

type sessionStoreWrap struct {
	sessions.Store
}

func (s *sessionStoreWrap) Options(opt session.Options) {}

func MiddlewareSessionWithConfig(name string, config SessionConfig) echo.MiddlewareFunc {
	_store, err := store.New(config.DB, config.TransformToSessionConfig(), []byte(Cfg.Session.SecretKey))

	if err != nil {
		panic(err)
	}
	return session.Sessions(name, &sessionStoreWrap{_store})
}

func MiddlewareInitSession() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			uri := ctx.Request().URI()

			if ctx.Get(session.DefaultKey) == nil {
				// TODO: reset all session?

				logrus.WithError(fmt.Errorf("not found session")).WithField("uri", uri).Fatal("current session is null")
				return nil
			}

			internalSession, ok := ctx.Get(session.DefaultKey).(session.Session)

			if !ok {
				// TODO: reset all session?

				logrus.WithError(fmt.Errorf("not valid type session")).WithField("uri", uri).Fatalf("current session is not Session, got %T", internalSession)
				return nil
			}

			if nil == internalSession {
				// should not happen

				return ctx.NoContent(http.StatusInternalServerError)
			}

			_session := NewSession(internalSession)

			ctx.Set(ctxKeyNameSession, _session)

			if _session.IsNew() {
				_session.Save()
				logrus.Debug("saved new session")
			}

			return h(ctx)
		}
	}
}
