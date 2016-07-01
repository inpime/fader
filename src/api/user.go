package api

// import (
// 	"github.com/inpime/dbox"
// 	"store"
// 	"strings"
// 	"utils"
// )

// var (
// 	UserEntityType dbox.EntityType = "user"

// 	UserLicensePropName = "licenses"

// 	GuestLicense = "guest"
// 	UserLicense  = "user"
// 	AdminLicense = "admin"
// )

// func FileAsUser(file *store.File) *User {
// 	return &User{
// 		File: *file,
// 	}
// }

// type User struct {
// 	store.File
// }

// func (u User) Type() dbox.EntityType {
// 	return UserEntityType
// }

// func (u User) Licenses() []string {

// 	return utils.Map(u.MapData()).Strings(UserLicensePropName)
// }

// func (u User) AddLicense(str string) User {
// 	// TODO: check valid license

// 	str = toLower(str)

// 	newLicenses := utils.NewA(utils.Map(u.MapData()).Strings(UserLicensePropName)).Add(str).(utils.AStrings).Array()
// 	utils.M(u.MapData()).Set(UserLicensePropName, newLicenses)

// 	return u
// }

// func (u User) DeleteLicense(str string) User {
// 	// TODO: check valid license

// 	str = toLower(str)

// 	newLicenses := utils.NewA(utils.Map(u.MapData()).Strings(UserLicensePropName)).Delete(str).(utils.AStrings).Array()
// 	utils.M(u.MapData()).Set(UserLicensePropName, newLicenses)

// 	return u
// }

// func (u User) HasLicense(str string) bool {
// 	// TODO: check valid license

// 	str = toLower(str)
// 	return utils.NewA(utils.Map(u.MapData()).Strings(UserLicensePropName)).Include(str)
// }

// //

// func toLower(str string) string {
// 	return strings.ToLower(str)
// }
