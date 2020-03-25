package core

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Hasher interface {
	HashPassword()
}

func (u *User) HashPassword() {
	//TODO: Handle error
	//TODO: Cost u daha iyi seçebilir misin
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hashedPassword)
}
