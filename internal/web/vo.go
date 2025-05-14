package web

type CreateBizReq struct {
	Name string
}

type BusinessConfig struct {
	ID        int64  `json:"id,omitzero"`
	OwnerID   int64  `json:"ownerId,omitzero"`
	OwnerType string `json:"ownerType,omitzero"`
	Name      string `json:"name,omitzero"`
	RateLimit int32  `json:"rateLimit,omitzero"`
	Token     string `json:"token,omitzero"`
}

type BusinessConfigReq struct {
	BizID          int64          `json:"bizId,omitzero"`
	BusinessConfig BusinessConfig `json:"businessConfig,omitzero"`
}

type ListReq struct {
	BizID  int64 `json:"bizId,omitzero"`
	Offset int   `json:"offset,omitzero"`
	Limit  int   `json:"limit,omitzero"`
}

type ListResp[T any] struct {
	Total int32 `json:"total,omitzero"`
	Rows  []T   `json:"Rows,omitzero"`
}

// Resource 资源
type Resource struct {
	ID          int64  `json:"id,omitzero"`
	BizID       int64  `json:"bizId,omitzero"`
	Type        string `json:"type,omitzero"`
	Key         string `json:"key,omitzero"`
	Name        string `json:"name,omitzero"`
	Description string `json:"description,omitzero"`
	Metadata    string `json:"metadata,omitzero"`
}
type ResourceReq struct {
	BizID    int64    `json:"bizId,omitzero"`
	Resource Resource `json:"resource,omitzero"`
}

type Permission struct {
	ID           int64  `json:"id,omitzero"`
	BizID        int64  `json:"bizID,omitzero"`
	Name         string `json:"name,omitzero"`
	Description  string `json:"description,omitzero"`
	ResourceID   int64  `json:"resourceId,omitzero"`
	ResourceType string `json:"resourceType,omitzero"`
	ResourceKey  string `json:"resourceKey,omitzero"`
	Action       string `json:"action,omitzero"`
	Metadata     string `json:"metadata,omitzero"`
}

type PermissionReq struct {
	BizID      int64      `json:"bizId,omitzero"`
	Permission Permission `json:"permission,omitzero"`
}

type Role struct {
	ID          int64  `json:"id,omitzero"`
	BizID       int64  `json:"bizID,omitzero"`
	Name        string `json:"name,omitzero"`
	Description string `json:"description,omitzero"`
	Metadata    string `json:"metadata,omitzero"`
}

type RoleReq struct {
	BizID int64 `json:"bizId,omitzero"`
	Role  Role  `json:"role,omitzero"`
}

type RoleInclusion struct {
	ID            int64 `json:"id,omitzero"`
	BizID         int64 `json:"bizID,omitzero"`
	IncludingRole Role  `json:"includingRole"`
	IncludedRole  Role  `json:"includedRole"`
}

type RoleInclusionReq struct {
	BizID         int64         `json:"bizId,omitzero"`
	RoleInclusion RoleInclusion `json:"roleInclusion,omitzero"`
}

type RolePermission struct {
	ID         int64      `json:"id,omitzero"`
	BizID      int64      `json:"bizID,omitzero"`
	Role       Role       `json:"role,omitzero"`
	Permission Permission `json:"permission,omitzero"`
}

type RolePermissionReq struct {
	BizID          int64          `json:"bizId,omitzero"`
	RolePermission RolePermission `json:"rolePermission,omitzero"`
}

type UserRole struct {
	ID        int64 `json:"id,omitzero"`
	BizID     int64 `json:"bizID,omitzero"`
	UserID    int64 `json:"userID,omitzero"`
	Role      Role  `json:"role"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

type UserRoleReq struct {
	BizID    int64    `json:"bizId,omitzero"`
	UserRole UserRole `json:"userRole,omitzero"`
}

type UserPermission struct {
	ID         int64      `json:"id,omitzero"`
	BizID      int64      `json:"bizID,omitzero"`
	UserID     int64      `json:"userID,omitzero"`
	Permission Permission `json:"permission,omitzero"`
	StartTime  int64      `json:"startTime,omitzero"`
	EndTime    int64      `json:"endTime,omitzero"`
	Effect     string     `json:"effect,omitzero"`
}

type UserPermissionReq struct {
	BizID          int64          `json:"bizId,omitzero"`
	UserPermission UserPermission `json:"userPermission,omitzero"`
}

type CreateAccountRoleReq struct {
	BizID int64 `json:"bizId,omitzero"`
	Role  Role  `json:"role,omitzero"`
}

type GrantAccountRolePermissionReq struct {
	BizID int64 `json:"bizId,omitzero"`
	Role       Role       `json:"role,omitzero"`
	Permission Permission `json:"permission,omitzero"`
}

type GrantAccountRoleReq struct {
	BizID int64 `json:"bizId,omitzero"`
	Role       Role       `json:"role,omitzero"`
	Permission Permission `json:"permission,omitzero"`
}

type RevokeRolePermissionReq struct {
	BizID int64 `json:"bizId,omitzero"`
	ID int64 `json:"id"`
}

type GrantUserRoleReq struct {
	BizID int64 `json:"bizId,omitzero"`
	UserID int64 `json:"userId,omitzero"`
	Role   Role  `json:"role,omitzero"`
}

type RevokeUserRoleReq struct {
	BizID int64 `json:"bizId,omitzero"`
	ID int64 `json:"id"`
}
