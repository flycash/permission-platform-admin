package domain

import "fmt"

type SystemTableResource string

const (
	BusinessConfigTable SystemTableResource = "business_configs"
	ResourceTable       SystemTableResource = "resources"
	PermissionTable     SystemTableResource = "permissions"
	RoleTable           SystemTableResource = "roles"
	RoleInclusionTable  SystemTableResource = "role_inclusions"
	RolePermissionTable SystemTableResource = "role_permissions"
	UserRoleTable       SystemTableResource = "user_roles"
	UserPermissionTable SystemTableResource = "user_permissions"
)

func (rk SystemTableResource) Type() string {
	return "system_table"
}

func (rk SystemTableResource) KeyForSystemAdmin() string {
	return fmt.Sprintf("/admin/%s", rk)
}

func (rk SystemTableResource) KeyForBusinessAdmin(bizID int64) string {
	return fmt.Sprintf("/admin/%s/%d", rk, bizID)
}

func (rk SystemTableResource) String() string {
	return string(rk)
}

type AccountResource string

const (
	ManagerAccountResource AccountResource = "account"
)

func (a AccountResource) String() string {
	return string(a)
}

func (a AccountResource) Type() string {
	return "admin_account"
}

func (a AccountResource) KeyForBusinessAdmin(bizID int64) string {
	return fmt.Sprintf("/admin/account/%d", bizID)
}

type PermissionActionType string

const (
	PermissionActionWrite PermissionActionType = "write"
	PermissionActionRead  PermissionActionType = "read"
)

func (a PermissionActionType) String() string {
	return string(a)
}
