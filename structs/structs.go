package structs

import "time"

// User ...
type User struct {
	ID        string
	Username  string
	Session   []Session
	LastLogin time.Time
	Files     []File
}

// File ...
type File struct {
	ID      string
	name    string
	size    int64
	content []byte
}

// Session ...
type Session struct {
	ID string
	IP string
}
