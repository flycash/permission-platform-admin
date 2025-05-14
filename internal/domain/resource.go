package domain

type Resource struct {
	ID          int64
	BizID       int64
	Type        string
	Key         string
	Name        string
	Description string
	Metadata    string
}
