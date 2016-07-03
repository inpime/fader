package session

import (
	"api/config"
	"github.com/Sirupsen/logrus"
	"github.com/gebv/echo-session"
	"github.com/inpime/dbox"
	"store"
)

var (
	SessionUserIDKey = addonName + ".user_id"

	// SessionExistingFlag defining non-empty session
	SessionExistingFlag = "existing"

	// DefaultGuestSession
	DefaultGuestSession *store.File
)

func NewSession(s session.Session) *Session {

	return &Session{
		Session: s,
		file:    nil,
	}
}

type Session struct {
	session.Session
	file *store.File
}

func (s *Session) Logout() *Session {
	userId := s.UserID()

	logrus.WithField("_service", addonName).Debugf("logout %q", userId)

	if err := s.Clear().Save(); err != nil {
		logrus.WithField("_service", addonName).WithError(err).Errorf("error save session %q", userId)
	}

	s.file = nil
	return s
}

func (s *Session) SetUserID(id string) *Session {
	return s.Set(SessionUserIDKey, id)
}

func (s Session) UserID() string {
	userId, ok := s.Get(SessionUserIDKey).(string)

	if !ok {
		return ""
	}

	return userId
}

func (s Session) IsNew() bool {
	return s.Get(SessionExistingFlag) == nil
}

func (s *Session) LoadFrom(file *store.File) *Session {
	logrus.WithField("_service", addonName).Debugf("load from %q [%q]", file.Name(), file.ID())
	s.file = file
	s.SetUserID(file.ID())

	return s
}

func (s *Session) AuthFrom(file *store.File) error {
	s.LoadFrom(file)

	return s.Save()
}

func (s *Session) User() *User {
	var userId = s.UserID()

	if s.file == nil && len(userId) > 0 {
		// user is authorized

		_file, err := dbox.NewFileID(userId, dbox.MustStore(config.UsersStoreName))
		file := store.MustFile(_file)

		logrus.WithFields(logrus.Fields{
			"_service":             addonName,
			"from_session:user_id": userId,
			"_loaded:user_id":      file.ID(),
			"_loaded:name":         file.Name(),
		}).Debugf("session existing, load")

		if err != nil {
			logrus.WithError(err).WithField("user_id", userId).Error("load user session")
		} else {
			s.LoadFrom(file)
		}
	}

	if s.file == nil {
		if err := s.AuthFrom(DefaultGuestSession); err != nil {
			logrus.WithError(err).WithField("user_id", DefaultGuestSession.ID()).Error("save guest user session")
		}
	}

	return FileAsUser(s.file)
}

// ------------------
// Default session interface
// ------------------

func (s Session) Get(key interface{}) interface{} {

	return s.Session.Get(key)
}

func (s *Session) Set(key interface{}, val interface{}) *Session {
	s.Session.Set(key, val)
	return s
}

func (s *Session) Delete(key interface{}) *Session {
	s.Session.Delete(key)
	return s
}

func (s *Session) Clear() *Session {
	s.Session.Clear()
	return s
}

func (s *Session) AddFlash(value interface{}, vars ...string) *Session {
	s.Session.AddFlash(value, vars...)
	return s
}

func (s Session) Flashes(vars ...string) []interface{} {
	messages := s.Session.Flashes(vars...)
	if err := s.Save(); err != nil {
		logrus.WithFields(logrus.Fields{
			"_service": addonName,
		}).WithError(err).Error("save session after get the flash messages")
	}
	return messages
}

func (s *Session) Save() error {
	s.Set(SessionExistingFlag, true)

	return s.Session.Save()
}
