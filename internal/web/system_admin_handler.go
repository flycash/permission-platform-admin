package web

import (
	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/ecodeclub/ginx"
	"github.com/ecodeclub/ginx/session"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
)

type SystemAdminHandler struct {
	*BaseHandler
	logger *elog.Component
}

func NewSystemAdminHandler(handler *BaseHandler) *SystemAdminHandler {
	return &SystemAdminHandler{BaseHandler: handler, logger: elog.DefaultLogger}
}

//nolint:dupl // 忽略
func (h *SystemAdminHandler) PrivateRoutes(server *gin.Engine) {
	server.POST("/admin/resource/create", ginx.BS[ResourceReq](h.CreateResource))
	server.GET("/admin/resource/get", ginx.BS[ResourceReq](h.GetResource))
	server.GET("/admin/resource/list", ginx.BS[ListReq](h.ListResources))
	server.POST("/admin/resource/update", ginx.BS[ResourceReq](h.UpdateResource))
	server.POST("/admin/resource/delete", ginx.BS[ResourceReq](h.DeleteResource))

	server.POST("/admin/permission/create", ginx.BS[PermissionReq](h.CreatePermission))
	server.GET("/admin/permission/get", ginx.BS[PermissionReq](h.GetPermission))
	server.GET("/admin/permission/list", ginx.BS[ListReq](h.ListPermissions))
	server.POST("/admin/permission/update", ginx.BS[PermissionReq](h.UpdatePermission))
	server.POST("/admin/permission/delete", ginx.BS[PermissionReq](h.DeletePermission))

	server.POST("/admin/role/create", ginx.BS[RoleReq](h.CreateRole))
	server.GET("/admin/role/get", ginx.BS[RoleReq](h.GetRole))
	server.GET("/admin/role/list", ginx.BS[ListReq](h.ListRoles))
	server.POST("/admin/role/update", ginx.BS[RoleReq](h.UpdateRole))
	server.POST("/admin/role/delete", ginx.BS[RoleReq](h.DeleteRole))

	server.POST("/admin/role-inclusion/create", ginx.BS[RoleInclusionReq](h.CreateRoleInclusion))
	server.GET("/admin/role-inclusion/get", ginx.BS[RoleInclusionReq](h.GetRoleInclusion))
	server.GET("/admin/role-inclusion/list", ginx.BS[ListReq](h.ListRoleInclusions))
	server.POST("/admin/role-inclusion/delete", ginx.BS[RoleInclusionReq](h.DeleteRoleInclusion))

	server.POST("/admin/role/grant_permission", ginx.BS[RolePermissionReq](h.GrantRolePermission))
	server.GET("/admin/role/list_permission", ginx.BS[ListReq](h.ListRolePermissions))
	server.POST("/admin/role/revoke_permission", ginx.BS[RolePermissionReq](h.RevokeRolePermission))

	server.POST("/admin/user/grant_role", ginx.BS(h.GrantUserRole))
	server.GET("/admin/user/list_role", ginx.BS[ListReq](h.ListUserRoles))
	server.POST("/admin/user/revoke_role", ginx.BS(h.RevokeUserRole))

	server.POST("/admin/user/grant_permission", ginx.BS[UserPermissionReq](h.GrantUserPermission))
	server.GET("/admin/user/list_permission", ginx.BS[ListReq](h.ListUserPermissions))
	server.POST("/admin/user/revoke_permission", ginx.BS[UserPermissionReq](h.RevokeUserPermission))
}

func (h *SystemAdminHandler) createCheckPermissionRequest(bizID, uid int64, resource domain.SystemTableResource, action domain.PermissionActionType) *permissionv1.CheckPermissionRequest {
	return &permissionv1.CheckPermissionRequest{
		Uid: uid,
		Permission: &permissionv1.Permission{
			BizId:        bizID,
			ResourceType: resource.Type(),
			ResourceKey:  resource.KeyForSystemAdmin(),
			Actions:      []string{action.String()},
		},
	}
}

// BusinessConfig

func (h *SystemAdminHandler) CreateBusinessConfig(ctx *ginx.Context, req BusinessConfigReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.BusinessConfigTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createBusinessConfig(systemAdminCtx, req)
}

func (h *SystemAdminHandler) GetBusinessConfig(ctx *ginx.Context, req BusinessConfigReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.BusinessConfigTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getBusinessConfig(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListBusinessConfigs(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.BusinessConfigTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listBusinessConfigs(systemAdminCtx, req)
}

func (h *SystemAdminHandler) UpdateBusinessConfig(ctx *ginx.Context, req BusinessConfigReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.BusinessConfigTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updateBusinessConfig(systemAdminCtx, req)
}

func (h *SystemAdminHandler) DeleteBusinessConfig(ctx *ginx.Context, req BusinessConfigReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.BusinessConfigTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteBusinessConfig(systemAdminCtx, req)
}

// Resource

func (h *SystemAdminHandler) CreateResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	h.logger.Info("invoked = 1")
	systemAdminCtx := h.systemAdminCtx(ctx)
	h.logger.Info("invoked = 2")
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionWrite))
	if err != nil {
		h.logger.Info("invoked = 3")
		return ginx.Result{}, err
	}
	h.logger.Info("invoked = 4")
	return h.createResource(systemAdminCtx, req)
}

func (h *SystemAdminHandler) GetResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getResource(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListResources(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listResources(systemAdminCtx, req)
}

func (h *SystemAdminHandler) UpdateResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updateResource(systemAdminCtx, req)
}

func (h *SystemAdminHandler) DeleteResource(ctx *ginx.Context, req ResourceReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.ResourceTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteResource(systemAdminCtx, req)
}

// Permission

func (h *SystemAdminHandler) CreatePermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createPermission(systemAdminCtx, req)
}

func (h *SystemAdminHandler) GetPermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getPermission(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListPermissions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listPermissions(systemAdminCtx, req)
}

func (h *SystemAdminHandler) UpdatePermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updatePermission(systemAdminCtx, req)
}

func (h *SystemAdminHandler) DeletePermission(ctx *ginx.Context, req PermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.PermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deletePermission(systemAdminCtx, req)
}

// Role

func (h *SystemAdminHandler) CreateRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createRole(systemAdminCtx, req)
}

func (h *SystemAdminHandler) GetRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getRole(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListRoles(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listRoles(systemAdminCtx, req)
}

func (h *SystemAdminHandler) UpdateRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.updateRole(systemAdminCtx, req)
}

func (h *SystemAdminHandler) DeleteRole(ctx *ginx.Context, req RoleReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteRole(systemAdminCtx, req)
}

// RoleInclusion

func (h *SystemAdminHandler) CreateRoleInclusion(ctx *ginx.Context, req RoleInclusionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.createRoleInclusion(systemAdminCtx, req)
}

func (h *SystemAdminHandler) GetRoleInclusion(ctx *ginx.Context, req RoleInclusionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.getRoleInclusion(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListRoleInclusions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listRoleInclusions(systemAdminCtx, req)
}

func (h *SystemAdminHandler) DeleteRoleInclusion(ctx *ginx.Context, req RoleInclusionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RoleInclusionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.deleteRoleInclusion(systemAdminCtx, req)
}

// RolePermission

func (h *SystemAdminHandler) GrantRolePermission(ctx *ginx.Context, req RolePermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RolePermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.grantRolePermission(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListRolePermissions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RolePermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listRolePermissions(systemAdminCtx, req)
}

func (h *SystemAdminHandler) RevokeRolePermission(ctx *ginx.Context, req RolePermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.RolePermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.revokeRolePermission(systemAdminCtx, req)
}

// UserRole

func (h *SystemAdminHandler) GrantUserRole(ctx *ginx.Context, req UserRoleReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserRoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.grantUserRole(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListUserRoles(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserRoleTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listUserRoles(systemAdminCtx, req)
}

func (h *SystemAdminHandler) RevokeUserRole(ctx *ginx.Context, req UserRoleReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserRoleTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.revokeUserRole(systemAdminCtx, req)
}

// UserPermission

func (h *SystemAdminHandler) GrantUserPermission(ctx *ginx.Context, req UserPermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserPermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.grantUserPermission(systemAdminCtx, req)
}

func (h *SystemAdminHandler) ListUserPermissions(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserPermissionTable, domain.PermissionActionRead))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.listUserPermissions(systemAdminCtx, req)
}

func (h *SystemAdminHandler) RevokeUserPermission(ctx *ginx.Context, req UserPermissionReq, sess session.Session) (ginx.Result, error) {
	systemAdminCtx := h.systemAdminCtx(ctx)
	err := h.checkPermission(systemAdminCtx, h.createCheckPermissionRequest(req.BizID, sess.Claims().Uid, domain.UserPermissionTable, domain.PermissionActionWrite))
	if err != nil {
		return ginx.Result{}, err
	}
	return h.revokeUserPermission(systemAdminCtx, req)
}
