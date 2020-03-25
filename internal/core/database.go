package core

import "database/sql"

type Database struct {
	DB *sql.DB
}

func (db *Database) AddUser(u User) error {
	statement, err := db.DB.Prepare("INSERT INTO users values (?, ?);")

	statement.Exec(u.Username, u.Password)
	return err
}
