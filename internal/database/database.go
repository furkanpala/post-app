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

func CheckUserExists(username string) (bool, error) {
	found := false

	rows, err := db.Query("SELECT username FROM users")

	if err != nil {
		return found, err
	}

	var u string

	for rows.Next() && !found {
		err := rows.Scan(&u)

		if err != nil {
			return found, err
		}

		if u == username {
			found = true
		}
	}

	return found, nil
}

func AddUser(user *core.User) error {
	statement, err := db.Prepare("INSERT INTO users(username,password) values (?, ?)")
	if err != nil {
		return err
	}
	statement.Exec(user.Username, user.Password)

	return nil
}
