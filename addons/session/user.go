package session

import (
	"strings"

	"github.com/inpime/dbox"
	"github.com/inpime/fader/store"
	"github.com/inpime/sdata"
)

var (
	UserEntityType dbox.EntityType = "user"

	UserLicensePropName = "licenses"

	GuestLicense = "guest"
	UserLicense  = "user"
	AdminLicense = "admin"
)

func FileAsUser(file *store.File) *User {
	return &User{
		File: *file,
	}
}

type User struct {
	store.File
}

func (u User) Type() dbox.EntityType {
	return UserEntityType
}

func (u *User) Licenses() *sdata.Array {

	return u.MMapData().A(UserLicensePropName)
}

func (u *User) AddLicense(str string) *User {
	// TODO: check valid license

	str = toLower(str)

	u.Licenses().Add(str)

	return u
}

func (u *User) DeleteLicense(str string) *User {
	// TODO: check valid license

	str = toLower(str)

	u.Licenses().Remove(str)

	return u
}

func (u User) HasLicense(str string) bool {
	// TODO: check valid license

	str = toLower(str)

	return u.Licenses().Includes(str)
}

//

func toLower(str string) string {
	return strings.ToLower(str)
}
