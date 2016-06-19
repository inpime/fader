package api

import (
	"github.com/gorilla/sessions"
	"store"
	"utils"
)

type gorillaSessionIface interface {
	Session() *sessions.Session
}

func (s Session) IsGuest() bool {

	return s.User().HasLicense(GuestLicense)
}

func (s Session) IsUser() bool {
	return s.User().HasLicense(UserLicense)
}

func (s Session) IsAuth() bool {
	return !s.HasLicense(GuestLicense)
}

func (s Session) HasLicense(name string) bool {
	return s.User().HasLicense(name)
}

//HasOneLicense has at least one license
func (s Session) HasOneLicense(names []string) bool {

	for _, name := range names {
		if s.User().HasLicense(name) {
			return true
		}
	}

	return false
}

func (s Session) Props() utils.M {
	return utils.Map(s.Session.(gorillaSessionIface).Session().Values)
}

// Signin the user authorization via key and password
func (s *Session) Signin(primaryKey, password string) error {

	file, err := store.LoadOrNewFile(UsersBucketName, toLower(primaryKey))

	if err != nil {
		return err
	}

	// TODO: encrypt a password and check password

	return s.Auth(file)
}
