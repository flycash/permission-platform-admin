package domain

const (
	DefaultAccountRoleType  = "admin_account"
	DefaultBusinessRoleType = "business_role"
)

// Role 角色
type Role struct {
	ID          int64
	BizID       int64
	Type        string
	Name        string
	Description string
	Metadata    string
}
