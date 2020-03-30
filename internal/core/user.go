package core

import "golang.org/x/crypto/bcrypt"

// User struct is a container to store user's username and password
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HashPassword function hashes the user's password with 10 salt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare function checks the plain text password against hashed password
func (u *User) Compare(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
