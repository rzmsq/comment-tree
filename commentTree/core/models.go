package core

type User struct {
	ID       int64
	Username string
}

type Comment struct {
	ID       int64
	UserID   int64
	ParentID int64
	Text     string
}
