package database

import (
	"database/sql"

	"github.com/furkanpala/post-app/internal/core"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./post-app.db")

	return err
}

func CloseDatabase() error {
	err := db.Close()
	return err
}

func CreateUsersTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users" (
		"username"	TEXT,
		"password"	TEXT NOT NULL,
		PRIMARY KEY("username")
	);`)
	statement.Exec()

	return err
}

func CreatePostsTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "posts" (
		"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"title"	TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		"sent_by"	TEXT NOT NULL,
		"date_added"	TEXT NOT NULL,
		FOREIGN KEY("sent_by") REFERENCES "users"("username")
	);`)
	statement.Exec()

	return err
}

func CreateLikesTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "likes" (
		"user_liked"	TEXT,
		"post_liked"	INTEGER,
		FOREIGN KEY("user_liked") REFERENCES "users"("username"),
		PRIMARY KEY("user_liked","post_liked"),
		FOREIGN KEY("post_liked") REFERENCES "posts"("id")
	);`)
	statement.Exec()

	return err
}

func CreateBlacklistTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "blacklist" (
		"refreshTokenString"	TEXT,
		PRIMARY KEY("refreshTokenString")
	);`)
	statement.Exec()

	return err
}

func FindUser(username string) (*core.User, error) {
	var user *core.User

	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		return user, err
	}

	var dbUsername string
	var dbPassword string

	for rows.Next() && user == nil {
		err := rows.Scan(&dbUsername, &dbPassword)

		if err != nil {
			return user, err
		}

		if username == dbUsername {
			user = &core.User{
				Username: dbUsername,
				Password: dbPassword,
			}
		}
	}

	return user, nil
}

func AddUser(user *core.User) error {
	statement, err := db.Prepare("INSERT INTO users(username,password) values (?, ?)")
	if err != nil {
		return err
	}
	statement.Exec(user.Username, user.Password)

	return nil
}

func BlacklistToken(token string) error {
	statement, err := db.Prepare("INSERT INTO blaclist values (?)")
	if err != nil {
		return err
	}
	statement.Exec(token)

	return nil
}

func FindToken(token string) (bool, error) {
	found := false
	rows, err := db.Query("SELECT * FROM blacklist")
	if err != nil {
		return false, err
	}

	var dbToken string
	for rows.Next() && !found {
		if err := rows.Scan(&dbToken); err != nil {
			return false, err
		}
		if dbToken == token {
			found = true
		}
	}

	return found, nil
}
