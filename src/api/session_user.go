package api

// import (
// 	// "fmt"
// 	"github.com/Sirupsen/logrus"
// 	"github.com/gebv/echo-session"
// 	"github.com/inpime/dbox"
// 	"store"
// 	"utils"
// )

// var (
// 	SessionUserIDKey = "user_id"

// 	// SessionExistingFlag defining non-empty session
// 	SessionExistingFlag = "existing"

// 	// SessionExistingFlag session properties
// 	SessionPropsKey = "ses_props"

// 	// DefaultGuestSession
// 	DefaultGuestSession *store.File
// )

// func NewSession(s session.Session) *Session {

// 	return &Session{
// 		Session:  s,
// 		userFile: nil,
// 	}
// }

// type Session struct {
// 	session.Session

// 	userFile *User
// }

// // BucketFilesProvider провайдер файлов бакета
// func (s Session) Licenses() utils.A {

// 	return utils.NewA(s.User().Licenses())
// }

// func (s *Session) User() *User {
// 	var userId = s.UserID()

// 	if s.userFile == nil && len(userId) > 0 {
// 		// user is authorized

// 		_file, err := dbox.NewFileID(userId, dbox.MustStore(UsersStoreName))
// 		file := store.MustFile(_file)

// 		logrus.WithFields(logrus.Fields{
// 			"from_session:user_id": userId,
// 			"_loaded:user_id":      file.ID(),
// 			"_loaded:name":         file.Name(),
// 		}).Debugf("session existing, load")

// 		if err != nil {
// 			logrus.WithError(err).WithField("user_id", userId).Error("load user session")
// 		} else {
// 			s.Load(file)
// 		}
// 	}

// 	if s.userFile == nil {
// 		if err := s.Auth(DefaultGuestSession); err != nil {
// 			logrus.WithError(err).WithField("user_id", DefaultGuestSession.ID()).Error("save guest user session")
// 		}
// 	}

// 	return s.userFile
// }

// func (s Session) UserID() string {
// 	userId, ok := s.Get(SessionUserIDKey).(string)

// 	if !ok {
// 		return ""
// 	}

// 	return userId
// }

// func (s Session) IsNew() bool {
// 	return s.Get(SessionExistingFlag) == nil
// }

// func (s *Session) SetUserID(userId string) *Session {
// 	return s.Set(SessionUserIDKey, userId)
// }

// func (s *Session) Logout() *Session {
// 	userId := s.UserID()
// 	if err := s.Clear().Save(); err != nil {
// 		logrus.WithError(err).WithField("user_id", userId).Error("clear and save session")
// 	}
// 	s.userFile = nil

// 	return s
// }

// func (s *Session) Load(file *store.File) {
// 	logrus.WithFields(logrus.Fields{
// 		"_loaded:user_id": file.ID(),
// 		"_loaded:name":    file.Name(),
// 	}).Debugf("session auth")

// 	s.userFile = FileAsUser(file)
// 	s.SetUserID(file.ID())

// 	return
// }

// // Signin
// func (s *Session) Auth(file *store.File) error {
// 	s.Load(file)

// 	return s.Save()
// }

// // ------------------
// // Default session interface
// // ------------------

// func (s Session) Get(key interface{}) interface{} {

// 	return s.Session.Get(key)
// }

// func (s *Session) Set(key interface{}, val interface{}) *Session {
// 	s.Session.Set(key, val)
// 	return s
// }

// func (s *Session) Delete(key interface{}) *Session {
// 	s.Session.Delete(key)
// 	return s
// }

// func (s *Session) Clear() *Session {
// 	s.Session.Clear()
// 	return s
// }

// func (s *Session) AddFlash(value interface{}, vars ...string) *Session {
// 	s.Session.AddFlash(value, vars...)
// 	return s
// }

// func (s Session) Flashes(vars ...string) []interface{} {

// 	return s.Session.Flashes(vars...)
// }

// func (s *Session) Save() error {
// 	s.Set(SessionExistingFlag, true)

// 	return s.Session.Save()
// }
