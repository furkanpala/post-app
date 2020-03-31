package core

// Post struct is a container to store post's ID,title, content, user and date added.
type Post struct {
	ID      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
	User    string `json:"user,omitempty"`
	Date    int64  `json:"date,omitempty"`
}
