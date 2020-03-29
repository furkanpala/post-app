package core

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Compare(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
