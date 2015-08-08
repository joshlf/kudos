package user

import (
	"os/user"
)

// UserFromOSUser is meant to be used
// by testing code in other packages
func UserFromOSUser(u *user.User) User {
	return User{u}
}
