package database

import (
	"database/sql"
	"time"

	"github.com/furkanpala/post-app/internal/core"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB

// OpenDatabase function opens the database file in root directory of project
func OpenDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./post-app.db")

	return err
}

func CloseDatabase() error {
	err := db.Close()
	return err
}

// CreateUsersTable function creates users table if it does not exist
// users table stores the username and password
func CreateUsersTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "users" (
		"username"	TEXT,
		"password"	TEXT NOT NULL,
		PRIMARY KEY("username")
	);`)
	statement.Exec()

	return err
}

// CreatePostsTable function creates posts table if it does not exist
// posts table stores the context and writer of posts
func CreatePostsTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "posts" (
		"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"title"	TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		"sent_by"	TEXT NOT NULL,
		"date_added"	INTEGER NOT NULL,
		FOREIGN KEY("sent_by") REFERENCES "users"("username")
	);`)
	statement.Exec()

	return err
}

// CreateLikesTable function creates likes table if it does not exist
// likes table stores the posts and users who liked that post
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

// CreateBlacklistTable function creates blacklist table if it does not exist
// blacklist table is used for revoking refresh tokens
func CreateBlacklistTable() error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "blacklist" (
		"jti"	TEXT,
		"expiresAt"	INTEGER NOT NULL,
		PRIMARY KEY("jti")
	);`)
	statement.Exec()

	return err
}

// FindUser function searches database for a specific user.
// Returns a pointer to the core.User if it finds
// nil otherwise.
func FindUser(username string) (*core.User, error) {
	var user *core.User

	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	var dbUsername string
	var dbPassword string

	for rows.Next() && user == nil {
		err := rows.Scan(&dbUsername, &dbPassword)

		if err != nil {
			return nil, err
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

// AddUser function adds the username and password into users table in database.
func AddUser(user *core.User) error {
	statement, err := db.Prepare("INSERT INTO users(username,password) values (?, ?)")
	if err != nil {
		return err
	}
	statement.Exec(user.Username, user.Password)

	return nil
}

// BlacklistToken function expects a string, jti, and int64, expire time.
// Adds jti and expire time of the JWT into blacklist table in database.
func BlacklistToken(jti string, expiresAt int64) error {
	statement, err := db.Prepare("INSERT INTO blacklist values (?,?)")
	if err != nil {
		return err
	}
	statement.Exec(jti, expiresAt)

	return nil
}

// FindJTI function searches the given jti string in database.
func FindJTI(jti string) (bool, error) {
	found := false
	rows, err := db.Query("SELECT jti FROM blacklist")
	if err != nil {
		return false, err
	}

	var dbJti string
	for rows.Next() && !found {
		if err := rows.Scan(&dbJti); err != nil {
			return false, err
		}
		if dbJti == jti {
			found = true
		}
	}

	return found, nil
}

// CountPosts function returns the number of posts in database
func CountPosts() (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM posts")
	if err != nil {
		return 0, err
	}
	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return count, err
		}
	}

	return count, nil
}

// GetAllPosts returns all the posts inside database
func GetAllPosts() ([]core.Post, error) {
	count, err := CountPosts()
	if err != nil {
		return nil, err
	}

	posts := make([]core.Post, count)

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}

	var post core.Post

	for rows.Next() {
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.User, &post.Date); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func AddPost(post *core.Post) error {
	statement, err := db.Prepare("INSERT INTO posts(title,content,sent_by,date_added) values(?,?,?,?)")
	if err != nil {
		return err
	}

	statement.Exec(post.Title, post.Content, post.User, time.Now().Unix())
	return nil
}
