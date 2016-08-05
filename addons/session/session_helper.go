package session

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/store"
	"github.com/inpime/sdata"
)

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

type gorillaSessionIface interface {
	Session() *sessions.Session
}

func (s Session) Props() *sdata.Map {

	return sdata.NewMapFrom(s.Session.(gorillaSessionIface).Session().Values)
}

// Signin the user authorization via key and password
func (s *Session) Signin(primaryKey, password string) error {

	file, err := store.LoadOrNewFile(config.UsersBucketName, toLower(primaryKey))

	if err != nil {
		return err
	}

	if file.MMapData().String("pwd") != password {

		return fmt.Errorf("invalid credentials")
	}

	// TODO: encrypt a password and check password

	return s.AuthFrom(file)
}
