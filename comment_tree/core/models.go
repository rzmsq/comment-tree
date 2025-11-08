package core

type Comment struct {
	ID       *int64 `json:"id,omitempty" required:"false"`
	Username string `json:"username,omitempty" required:"false"`
	ParentID *int64 `json:"parent_id,omitempty" required:"false"`
	Text     string `json:"text" required:"true"`
}
