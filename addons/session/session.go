package session

import (
	"encoding/gob"

	"github.com/Sirupsen/logrus"
	session "github.com/echo-contrib/sessions"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/store"
	"github.com/inpime/fader/utils/sdata"
)

var (
	SessionUserIDKey = addonName + ".user_id"

	// SessionExistingFlag defining non-empty session
	SessionExistingFlag = "existing"

	// DefaultGuestSession
	DefaultGuestSession *store.File
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(sdata.NewStringMap())
}

func NewSession(s session.Session) *Session {

	_session := &Session{
		Session: s,
		file:    nil,
	}

	if s.Get(SessionExistingFlag) == nil {

		if err := _session.Save(); err != nil {
			logrus.WithFields(logrus.Fields{
				"_api": addonName,
			}).
				WithError(err).
				Error("save new session")
		}
	}

	return _session
}

type Session struct {
	Session session.Session
	file    *store.File
}

func (s *Session) Logout() *Session {
	userId := s.UserID()

	logrus.WithField("_api", addonName).Debugf("logout %q", userId)

	if err := s.Clear().Save(); err != nil {
		logrus.WithField("_api", addonName).WithError(err).Errorf("error save session %q", userId)
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

func (s *Session) IsNew() bool {
	return s.Get(SessionExistingFlag) == nil
}

func (s *Session) LoadFrom(file *store.File) *Session {
	logrus.WithField("_api", addonName).Debugf("load from %q [%q]", file.Name(), file.ID())
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
			"_api":                 addonName,
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

func (s *Session) Get(key interface{}) interface{} {

	return s.Session.Get(key)
}

// GetOnce get and remove field
func (s Session) GetOnce(key interface{}) interface{} {
	v := s.Session.Get(key)
	s.Session.Delete(key)

	return v
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

	return messages
}

// FirstFlash return the first of messages. All messages are cleared
func (s Session) FirstFlash(vars ...string) interface{} {
	messages := s.Flashes(vars...)

	if len(messages) > 0 {
		return messages[0]
	}

	return nil
}

func (s *Session) Save() error {
	if s.Get(SessionExistingFlag) == nil {
		s.Set(SessionExistingFlag, true)
	}

	return s.Session.Save()
}
