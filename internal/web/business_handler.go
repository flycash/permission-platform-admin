package web

import (
	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/ecodeclub/ginx"
	"github.com/ecodeclub/ginx/session"
	"github.com/gin-gonic/gin"
)

type BusinessHandler struct {
	BaseHandler
}

func (h *BusinessHandler) PrivateRoutes(server *gin.Engine) {
	server.POST("/resource/create", ginx.BS[ResourceReq](h.CreateResource))
	server.POST("/resource/get", ginx.BS[ResourceReq](h.GetResource))
	server.POST("/resource/list", ginx.BS[ListReq](h.ListResources))
	server.POST("/resource/update", ginx.BS[ResourceReq](h.UpdateResource))
	server.POST("/resource/delete", ginx.BS[ResourceReq](h.DeleteResource))

	server.POST("/permission/create", ginx.BS[PermissionReq](h.CreatePermission))
	server.POST("/permission/get", ginx.BS[PermissionReq](h.GetPermission))
	server.POST("/permission/list", ginx.BS[ListReq](h.ListPermissions))
	server.POST("/permission/update", ginx.BS[PermissionReq](h.UpdatePermission))
	server.POST("/permission/delete", ginx.BS[PermissionReq](h.DeletePermission))

	server.POST("/role/create", ginx.BS[RoleReq](h.CreateRole))
	server.POST("/role/get", ginx.BS[RoleReq](h.GetRole))
	server.POST("/role/list", ginx.BS[ListReq](h.ListRoles))
	server.POST("/role/update", ginx.BS[RoleReq](h.UpdateRole))
	server.POST("/role/delete", ginx.BS[RoleReq](h.DeleteRole))

	server.POST("/role-inclusion/create", ginx.BS[RoleInclusionReq](h.CreateRoleInclusion))
	server.POST("/role-inclusion/get", ginx.BS[RoleInclusionReq](h.GetRoleInclusion))
	server.POST("/role-inclusion/list", ginx.BS[ListReq](h.ListRoleInclusions))
	server.POST("/role-inclusion/delete", ginx.BS[RoleInclusionReq](h.DeleteRoleInclusion))

	server.POST("/role/grant_permission", ginx.BS[RolePermissionReq](h.GrantRolePermission))
	server.POST("/role/list_permission", ginx.BS[ListReq](h.ListRolePermissions))
	server.POST("/role/revoke_permission", ginx.BS[RolePermissionReq](h.RevokeRolePermission))

	server.POST("/user/grant_role", ginx.BS(h.GrantUserRole))
	server.POST("/user/list_role", ginx.BS[ListReq](h.ListUserRoles))
	server.POST("/user/revoke_role", ginx.BS(h.RevokeUserRole))

	server.POST("/user/grant_permission", ginx.BS[UserPermissionReq](h.GrantUserPermission))
	server.POST("/user/list_permission", ginx.BS[ListReq](h.ListUserPermissions))
	server.POST("/user/revoke_permission", ginx.BS[UserPermissionReq](h.RevokeUserPermission))
}

func (h *BusinessHandler) createCheckPermissionRequest(bizID, uid int64, resource domain.SystemTableResource, action domain.PermissionActionType) *permissionv1.CheckPermissionRequest {
	return &permissionv1.CheckPermissionRequest{
		Uid: uid,
		Permission: &permissionv1.Permission{
			BizId:        bizID,
			ResourceType: resource.Type(),
			ResourceKey:  resource.KeyForBusinessAdmin(bizID),
			Actions:      []string{action.String()},
		},
	}
}

// Resource

func (h *BusinessHandler) CreateResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createResource(businessAdminCtx, req)
}

func (h *BusinessHandler) GetResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getResource(businessAdminCtx, req)
}

func (h *BusinessHandler) ListResources(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listResources(businessAdminCtx, req)
}

func (h *BusinessHandler) UpdateResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updateResource(businessAdminCtx, req)
}

func (h *BusinessHandler) DeleteResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteResource(businessAdminCtx, req)
}

// Permission

func (h *BusinessHandler) CreatePermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createPermission(businessAdminCtx, req)
}

func (h *BusinessHandler) GetPermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getPermission(businessAdminCtx, req)
}

func (h *BusinessHandler) ListPermissions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listPermissions(businessAdminCtx, req)
}

func (h *BusinessHandler) UpdatePermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updatePermission(businessAdminCtx, req)
}

func (h *BusinessHandler) DeletePermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deletePermission(businessAdminCtx, req)
}

// Role

func (h *BusinessHandler) CreateRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createRole(businessAdminCtx, req)
}

func (h *BusinessHandler) GetRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getRole(businessAdminCtx, req)
}

func (h *BusinessHandler) ListRoles(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listRoles(businessAdminCtx, req)
}

func (h *BusinessHandler) UpdateRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updateRole(businessAdminCtx, req)
}

func (h *BusinessHandler) DeleteRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteRole(businessAdminCtx, req)
}

// RoleInclusion

func (h *BusinessHandler) CreateRoleInclusion(ctx *ginx.Context, req RoleInclusionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createRoleInclusion(businessAdminCtx, req)
}

func (h *BusinessHandler) GetRoleInclusion(ctx *ginx.Context, req RoleInclusionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getRoleInclusion(businessAdminCtx, req)
}

func (h *BusinessHandler) ListRoleInclusions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listRoleInclusions(businessAdminCtx, req)
}

func (h *BusinessHandler) DeleteRoleInclusion(ctx *ginx.Context, req RoleInclusionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteRoleInclusion(businessAdminCtx, req)
}

// RolePermission

func (h *BusinessHandler) GrantRolePermission(ctx *ginx.Context, req RolePermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RolePermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.grantRolePermission(businessAdminCtx, req)
}

func (h *BusinessHandler) ListRolePermissions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RolePermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listRolePermissions(businessAdminCtx, req)
}

func (h *BusinessHandler) RevokeRolePermission(ctx *ginx.Context, req RolePermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RolePermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.revokeRolePermission(businessAdminCtx, req)
}

// UserRole

func (h *BusinessHandler) GrantUserRole(ctx *ginx.Context, req UserRoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserRoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.grantUserRole(businessAdminCtx, req)
}

func (h *BusinessHandler) ListUserRoles(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserRoleTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listUserRoles(businessAdminCtx, req)
}

func (h *BusinessHandler) RevokeUserRole(ctx *ginx.Context, req UserRoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserRoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.revokeUserRole(businessAdminCtx, req)
}

// UserPermission

func (h *BusinessHandler) GrantUserPermission(ctx *ginx.Context, req UserPermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserPermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.grantUserPermission(businessAdminCtx, req)
}

func (h *BusinessHandler) ListUserPermissions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserPermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listUserPermissions(businessAdminCtx, req)
}

func (h *BusinessHandler) RevokeUserPermission(ctx *ginx.Context, req UserPermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkPermission(businessAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserPermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.revokeUserPermission(businessAdminCtx, req)
}
