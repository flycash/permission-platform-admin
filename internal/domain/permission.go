package domain

type Permission struct {
	ID          int64
	BizID       int64
	Name        string
	Description string
	Resource    Resource
	Action      string
	Metadata    string
}
