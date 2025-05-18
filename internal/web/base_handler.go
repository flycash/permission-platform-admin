package web

import (
	"context"
	"errors"

	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ginx"
	"google.golang.org/grpc/metadata"
)

type BaseHandler struct {
	rbacSvc       permissionv1.RBACServiceClient
	permissionSvc permissionv1.PermissionServiceClient
	adminToken    string
}

func NewBaseHandler(rbacSvc permissionv1.RBACServiceClient, permissionSvc permissionv1.PermissionServiceClient, adminToken string) *BaseHandler {
	return &BaseHandler{rbacSvc: rbacSvc, permissionSvc: permissionSvc, adminToken: adminToken}
}

func (h *BaseHandler) systemAdminCtx(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", h.adminToken)
}

func (h *BaseHandler) businessAdminCtx(ctx context.Context, bizID int64) (context.Context, error) {
	bizResp, err := h.rbacSvc.GetBusinessConfig(h.systemAdminCtx(ctx), &permissionv1.GetBusinessConfigRequest{
		Id: bizID,
	})
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, "Authorization", bizResp.Config.Token), nil
}

func (h *BaseHandler) checkPermission(ctx context.Context, checkPermissionRequest *permissionv1.CheckPermissionRequest) error {
	resp, err := h.permissionSvc.CheckPermission(ctx, checkPermissionRequest)
	if err != nil {
		return err
	}
	if !resp.Allowed {
		return errors.New("没有权限")
	}
	return nil
}

// BusinessConfig

func (h *BaseHandler) createBusinessConfig(ctx context.Context, req BusinessConfigReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.CreateBusinessConfig(ctx, &permissionv1.CreateBusinessConfigRequest{
		Config: h.toBusinessConfigPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Config.Id,
	}, nil
}

func (h *BaseHandler) toBusinessConfigPB(req BusinessConfigReq) *permissionv1.BusinessConfig {
	return &permissionv1.BusinessConfig{
		OwnerId:   req.BusinessConfig.OwnerID,
		OwnerType: req.BusinessConfig.OwnerType,
		Name:      req.BusinessConfig.Name,
		RateLimit: req.BusinessConfig.RateLimit,
	}
}

func (h *BaseHandler) getBusinessConfig(ctx context.Context, req BusinessConfigReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GetBusinessConfig(ctx, &permissionv1.GetBusinessConfigRequest{
		Id: req.BusinessConfig.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toBusinessConfigVO(resp.Config),
	}, nil
}

func (h *BaseHandler) toBusinessConfigVO(src *permissionv1.BusinessConfig) BusinessConfig {
	return BusinessConfig{
		ID:        src.Id,
		OwnerID:   src.OwnerId,
		OwnerType: src.OwnerType,
		Name:      src.Name,
		RateLimit: src.RateLimit,
		Token:     src.Token,
	}
}

func (h *BaseHandler) listBusinessConfigs(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListBusinessConfigs(ctx, &permissionv1.ListBusinessConfigsRequest{
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[BusinessConfig]{
			Rows: slice.Map(resp.Configs, func(_ int, src *permissionv1.BusinessConfig) BusinessConfig {
				return h.toBusinessConfigVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) updateBusinessConfig(ctx context.Context, req BusinessConfigReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.UpdateBusinessConfig(ctx, &permissionv1.UpdateBusinessConfigRequest{
		Config: h.toBusinessConfigPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

func (h *BaseHandler) deleteBusinessConfig(ctx context.Context, req BusinessConfigReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.DeleteBusinessConfig(ctx, &permissionv1.DeleteBusinessConfigRequest{
		Id: req.BusinessConfig.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// Resource

func (h *BaseHandler) createResource(ctx context.Context, req ResourceReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.CreateResource(ctx, &permissionv1.CreateResourceRequest{
		Resource: h.toResourcePB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Resource.Id,
	}, nil
}

func (h *BaseHandler) toResourcePB(req ResourceReq) *permissionv1.Resource {
	return &permissionv1.Resource{
		Id:          req.Resource.ID,
		BizId:       req.Resource.BizID,
		Type:        req.Resource.Type,
		Key:         req.Resource.Key,
		Name:        req.Resource.Name,
		Description: req.Resource.Description,
		Metadata:    req.Resource.Metadata,
	}
}

func (h *BaseHandler) getResource(ctx context.Context, req ResourceReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GetResource(ctx, &permissionv1.GetResourceRequest{
		Id: req.Resource.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toResourceVO(resp.Resource),
	}, nil
}

func (h *BaseHandler) toResourceVO(src *permissionv1.Resource) Resource {
	return Resource{
		ID:          src.Id,
		BizID:       src.BizId,
		Type:        src.Type,
		Key:         src.Key,
		Name:        src.Name,
		Description: src.Description,
		Metadata:    src.Metadata,
	}
}

func (h *BaseHandler) listResources(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListResources(ctx, &permissionv1.ListResourcesRequest{
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[Resource]{
			Rows: slice.Map(resp.Resources, func(_ int, src *permissionv1.Resource) Resource {
				return h.toResourceVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) updateResource(ctx context.Context, req ResourceReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.UpdateResource(ctx, &permissionv1.UpdateResourceRequest{
		Resource: h.toResourcePB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

func (h *BaseHandler) deleteResource(ctx context.Context, req ResourceReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.DeleteResource(ctx, &permissionv1.DeleteResourceRequest{
		Id: req.Resource.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// Permission

func (h *BaseHandler) createPermission(ctx context.Context, req PermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.CreatePermission(ctx, &permissionv1.CreatePermissionRequest{
		Permission: h.toPermissionPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Permission.Id,
	}, nil
}

func (h *BaseHandler) toPermissionPB(req PermissionReq) *permissionv1.Permission {
	return &permissionv1.Permission{
		Id:           req.Permission.ID,
		BizId:        req.Permission.BizID,
		Name:         req.Permission.Name,
		Description:  req.Permission.Description,
		ResourceId:   req.Permission.ResourceID,
		ResourceType: req.Permission.ResourceType,
		ResourceKey:  req.Permission.ResourceKey,
		Actions:      []string{req.Permission.Action},
		Metadata:     req.Permission.Metadata,
	}
}

func (h *BaseHandler) getPermission(ctx context.Context, req PermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GetPermission(ctx, &permissionv1.GetPermissionRequest{
		Id: req.Permission.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toPermissionVO(resp.Permission),
	}, nil
}

func (h *BaseHandler) toPermissionVO(src *permissionv1.Permission) Permission {
	const first = 0
	return Permission{
		ID:           src.Id,
		BizID:        src.BizId,
		Name:         src.Name,
		Description:  src.Description,
		ResourceID:   src.ResourceId,
		ResourceType: src.ResourceType,
		ResourceKey:  src.ResourceKey,
		Action:       src.Actions[first],
		Metadata:     src.Metadata,
	}
}

func (h *BaseHandler) listPermissions(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListPermissions(ctx, &permissionv1.ListPermissionsRequest{
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[Permission]{
			Rows: slice.Map(resp.Permissions, func(_ int, src *permissionv1.Permission) Permission {
				return h.toPermissionVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) updatePermission(ctx context.Context, req PermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.UpdatePermission(ctx, &permissionv1.UpdatePermissionRequest{
		Permission: h.toPermissionPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

func (h *BaseHandler) deletePermission(ctx context.Context, req PermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.DeleteResource(ctx, &permissionv1.DeleteResourceRequest{
		Id: req.Permission.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// Role

func (h *BaseHandler) createRole(ctx context.Context, req RoleReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.CreateRole(ctx, &permissionv1.CreateRoleRequest{
		Role: h.toRolePB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Role.Id,
	}, nil
}

func (h *BaseHandler) toRolePB(req RoleReq) *permissionv1.Role {
	return &permissionv1.Role{
		Id:          req.Role.ID,
		BizId:       req.Role.BizID,
		Type:        domain.DefaultBusinessRoleType,
		Name:        req.Role.Name,
		Description: req.Role.Description,
		Metadata:    req.Role.Metadata,
	}
}

func (h *BaseHandler) getRole(ctx context.Context, req RoleReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GetRole(ctx, &permissionv1.GetRoleRequest{
		Id: req.Role.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toRoleVO(resp.Role),
	}, nil
}

func (h *BaseHandler) toRoleVO(src *permissionv1.Role) Role {
	return Role{
		ID:          src.Id,
		BizID:       src.BizId,
		Name:        src.Name,
		Description: src.Description,
		Metadata:    src.Metadata,
	}
}

func (h *BaseHandler) listRoles(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListRoles(ctx, &permissionv1.ListRolesRequest{
		Type:   domain.DefaultBusinessRoleType,
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[Role]{
			Rows: slice.Map(resp.Roles, func(_ int, src *permissionv1.Role) Role {
				return h.toRoleVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) updateRole(ctx context.Context, req RoleReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.UpdateRole(ctx, &permissionv1.UpdateRoleRequest{
		Role: h.toRolePB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

func (h *BaseHandler) deleteRole(ctx context.Context, req RoleReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.DeleteResource(ctx, &permissionv1.DeleteResourceRequest{
		Id: req.Role.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// RoleInclusion

func (h *BaseHandler) createRoleInclusion(ctx context.Context, req RoleInclusionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.CreateRoleInclusion(ctx, &permissionv1.CreateRoleInclusionRequest{
		RoleInclusion: h.toRoleInclusionPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.RoleInclusion.Id,
	}, nil
}

func (h *BaseHandler) toRoleInclusionPB(req RoleInclusionReq) *permissionv1.RoleInclusion {
	return &permissionv1.RoleInclusion{
		Id:                req.RoleInclusion.ID,
		BizId:             req.RoleInclusion.BizID,
		IncludingRoleId:   req.RoleInclusion.IncludingRole.ID,
		IncludingRoleType: domain.DefaultBusinessRoleType,
		IncludingRoleName: req.RoleInclusion.IncludingRole.Name,
		IncludedRoleId:    req.RoleInclusion.IncludedRole.ID,
		IncludedRoleType:  domain.DefaultBusinessRoleType,
		IncludedRoleName:  req.RoleInclusion.IncludedRole.Name,
	}
}

func (h *BaseHandler) getRoleInclusion(ctx context.Context, req RoleInclusionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GetRoleInclusion(ctx, &permissionv1.GetRoleInclusionRequest{
		Id: req.RoleInclusion.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toRoleInclusionVO(resp.RoleInclusion),
	}, nil
}

func (h *BaseHandler) toRoleInclusionVO(src *permissionv1.RoleInclusion) RoleInclusion {
	return RoleInclusion{
		ID:    src.Id,
		BizID: src.BizId,
		IncludingRole: Role{
			ID:   src.IncludingRoleId,
			Name: src.IncludingRoleName,
		},
		IncludedRole: Role{
			ID:   src.IncludedRoleId,
			Name: src.IncludedRoleName,
		},
	}
}

func (h *BaseHandler) listRoleInclusions(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListRoleInclusions(ctx, &permissionv1.ListRoleInclusionsRequest{
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[RoleInclusion]{
			Rows: slice.Map(resp.RoleInclusions, func(_ int, src *permissionv1.RoleInclusion) RoleInclusion {
				return h.toRoleInclusionVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) deleteRoleInclusion(ctx context.Context, req RoleInclusionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.DeleteRoleInclusion(ctx, &permissionv1.DeleteRoleInclusionRequest{
		Id: req.RoleInclusion.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// RolePermission

func (h *BaseHandler) grantRolePermission(ctx context.Context, req RolePermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GrantRolePermission(ctx, &permissionv1.GrantRolePermissionRequest{
		RolePermission: h.toRolePermissionPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toRolePermissionVO(resp.RolePermission),
	}, nil
}

func (h *BaseHandler) toRolePermissionPB(req RolePermissionReq) *permissionv1.RolePermission {
	return &permissionv1.RolePermission{
		Id:               req.RolePermission.ID,
		BizId:            req.RolePermission.BizID,
		RoleId:           req.RolePermission.Role.ID,
		PermissionId:     req.RolePermission.Permission.ID,
		RoleName:         req.RolePermission.Role.Name,
		RoleType:         domain.DefaultBusinessRoleType,
		ResourceType:     req.RolePermission.Permission.ResourceType,
		ResourceKey:      req.RolePermission.Permission.ResourceKey,
		PermissionAction: req.RolePermission.Permission.Action,
	}
}

func (h *BaseHandler) toRolePermissionVO(src *permissionv1.RolePermission) RolePermission {
	return RolePermission{
		ID:    src.Id,
		BizID: src.BizId,
		Role: Role{
			ID:   src.RoleId,
			Name: src.RoleName,
		},
		Permission: Permission{
			ID:           src.PermissionId,
			ResourceType: src.ResourceType,
			ResourceKey:  src.ResourceKey,
			Action:       src.PermissionAction,
		},
	}
}

func (h *BaseHandler) listRolePermissions(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListRolePermissions(ctx, &permissionv1.ListRolePermissionsRequest{
		BizId: req.BizID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[RolePermission]{
			Rows: slice.Map(resp.RolePermissions, func(_ int, src *permissionv1.RolePermission) RolePermission {
				return h.toRolePermissionVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) revokeRolePermission(ctx context.Context, req RolePermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.RevokeRolePermission(ctx, &permissionv1.RevokeRolePermissionRequest{
		Id: req.RolePermission.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// UserRole

func (h *BaseHandler) grantUserRole(ctx context.Context, req UserRoleReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GrantUserRole(ctx, &permissionv1.GrantUserRoleRequest{
		UserRole: h.toUserRolePB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toUserRoleVO(resp.UserRole),
	}, nil
}

func (h *BaseHandler) toUserRolePB(req UserRoleReq) *permissionv1.UserRole {
	return &permissionv1.UserRole{
		Id:        req.UserRole.ID,
		BizId:     req.UserRole.BizID,
		UserId:    req.UserRole.UserID,
		RoleId:    req.UserRole.Role.ID,
		RoleName:  req.UserRole.Role.Name,
		RoleType:  domain.DefaultBusinessRoleType,
		StartTime: req.UserRole.StartTime,
		EndTime:   req.UserRole.EndTime,
	}
}

func (h *BaseHandler) toUserRoleVO(src *permissionv1.UserRole) UserRole {
	return UserRole{
		ID:     src.Id,
		BizID:  src.BizId,
		UserID: src.UserId,
		Role: Role{
			ID:   src.RoleId,
			Name: src.RoleName,
		},
		StartTime: src.StartTime,
		EndTime:   src.EndTime,
	}
}

func (h *BaseHandler) listUserRoles(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListUserRoles(ctx, &permissionv1.ListUserRolesRequest{
		BizId: req.BizID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[UserRole]{
			Rows: slice.Map(resp.UserRoles, func(_ int, src *permissionv1.UserRole) UserRole {
				return h.toUserRoleVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) revokeUserRole(ctx context.Context, req UserRoleReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.RevokeUserRole(ctx, &permissionv1.RevokeUserRoleRequest{
		Id: req.UserRole.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

// UserPermission

func (h *BaseHandler) grantUserPermission(ctx context.Context, req UserPermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.GrantUserPermission(ctx, &permissionv1.GrantUserPermissionRequest{
		UserPermission: h.toUserPermissionPB(req),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toUserPermissionVO(resp.UserPermission),
	}, nil
}

func (h *BaseHandler) toUserPermissionPB(req UserPermissionReq) *permissionv1.UserPermission {
	return &permissionv1.UserPermission{
		Id:               req.UserPermission.ID,
		BizId:            req.UserPermission.BizID,
		UserId:           req.UserPermission.UserID,
		PermissionId:     req.UserPermission.Permission.ID,
		PermissionName:   req.UserPermission.Permission.Name,
		ResourceType:     req.UserPermission.Permission.ResourceType,
		ResourceKey:      req.UserPermission.Permission.ResourceKey,
		PermissionAction: req.UserPermission.Permission.Action,
		StartTime:        req.UserPermission.StartTime,
		EndTime:          req.UserPermission.EndTime,
		Effect:           req.UserPermission.Effect,
	}
}

func (h *BaseHandler) toUserPermissionVO(src *permissionv1.UserPermission) UserPermission {
	return UserPermission{
		ID:     src.Id,
		BizID:  src.BizId,
		UserID: src.UserId,
		Permission: Permission{
			ID:     src.PermissionId,
			Name:   src.PermissionName,
			Action: src.PermissionAction,
		},
		StartTime: src.StartTime,
		EndTime:   src.EndTime,
	}
}

func (h *BaseHandler) listUserPermissions(ctx context.Context, req ListReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.ListUserPermissions(ctx, &permissionv1.ListUserPermissionsRequest{
		Offset: int32(req.Offset),
		Limit:  int32(req.Limit),
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: ListResp[UserPermission]{
			Rows: slice.Map(resp.UserPermissions, func(_ int, src *permissionv1.UserPermission) UserPermission {
				return h.toUserPermissionVO(src)
			}),
		},
	}, nil
}

func (h *BaseHandler) revokeUserPermission(ctx context.Context, req UserPermissionReq) (ginx.Result, error) {
	resp, err := h.rbacSvc.RevokeUserPermission(ctx, &permissionv1.RevokeUserPermissionRequest{
		Id: req.UserPermission.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}
