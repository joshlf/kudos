package user

import (
	"os/user"
)

type User struct {
	u *user.User
}

func (u User) ID() string {
	return u.u.Uid
}

func Current() (User, error) {
	u, err := user.Current()
	return User{u}, err
}
