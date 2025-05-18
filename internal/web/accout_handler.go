package web

import (
	"context"

	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ginx"
	"github.com/ecodeclub/ginx/session"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	*BaseHandler
}

func NewAccountHandler(handler *BaseHandler) *AccountHandler {
	return &AccountHandler{BaseHandler: handler}
}

func (h *AccountHandler) PrivateRoutes(server *gin.Engine) {
	// 创建角色
	server.POST("/account/role/create", ginx.BS[CreateAccountRoleReq](h.CreateRole))
	// 展示角色
	server.GET("/account/role/list", ginx.BS[ListReq](h.ListRoles))
	// 赋予角色权限
	server.POST("/account/role/grant_permission", ginx.BS[GrantAccountRolePermissionReq](h.GrantRolePermission))
	// 撤销角色权限
	server.POST("/account/role/revoke_permission", ginx.BS[RevokeRolePermissionReq](h.RevokeRolePermission))
	// 赋予用户角色
	server.POST("/account/user/grant_role", ginx.BS[GrantUserRoleReq](h.GrantUserRole))
	// 撤销用户角色
	server.POST("/account/user/revoke_role", ginx.BS[RevokeUserRoleReq](h.RevokeUserRole))
}

func (h *AccountHandler) checkAccountPermission(ctx context.Context, bizID, userID int64) error {
	return h.checkPermission(ctx, &permissionv1.CheckPermissionRequest{
		Uid: userID,
		Permission: &permissionv1.Permission{
			ResourceType: domain.ManagerAccountResource.Type(),
			ResourceKey:  domain.ManagerAccountResource.KeyForBusinessAdmin(bizID),
			Actions:      []string{domain.PermissionActionRead.String(), domain.PermissionActionWrite.String()},
		},
	})
}

func (h *AccountHandler) CreateRole(ctx *ginx.Context, req CreateAccountRoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkAccountPermission(businessAdminCtx, req.BizID, sess.Claims().Uid)
	if err != nil {
		return ginx.Result{}, err
	}
	resp, err := h.rbacSvc.CreateRole(businessAdminCtx, &permissionv1.CreateRoleRequest{
		Role: &permissionv1.Role{
			Type:        domain.DefaultAccountRoleType,
			Name:        req.Role.Name,
			Description: req.Role.Description,
			Metadata:    req.Role.Metadata,
		},
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Role.Id,
	}, nil
}

func (h *AccountHandler) ListRoles(ctx *ginx.Context, req ListReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkAccountPermission(businessAdminCtx, req.BizID, sess.Claims().Uid)
	if err != nil {
		return ginx.Result{}, err
	}
	resp, err := h.rbacSvc.ListRoles(businessAdminCtx, &permissionv1.ListRolesRequest{
		Type:   domain.DefaultAccountRoleType,
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

func (h *AccountHandler) GrantRolePermission(ctx *ginx.Context, req GrantAccountRolePermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	if err := h.checkAccountPermission(businessAdminCtx, req.BizID, sess.Claims().Uid); err != nil {
		return ginx.Result{}, err
	}
	resp, err := h.rbacSvc.GrantRolePermission(businessAdminCtx, &permissionv1.GrantRolePermissionRequest{
		RolePermission: &permissionv1.RolePermission{
			RoleId:           req.Role.ID,
			PermissionId:     req.Permission.ID,
			RoleName:         req.Role.Name,
			RoleType:         domain.DefaultAccountRoleType,
			ResourceType:     req.Permission.ResourceType,
			ResourceKey:      req.Permission.ResourceKey,
			PermissionAction: req.Permission.Action,
		},
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toRolePermissionVO(resp.RolePermission),
	}, nil
}

func (h *AccountHandler) RevokeRolePermission(ctx *ginx.Context, req RevokeRolePermissionReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkAccountPermission(businessAdminCtx, req.BizID, sess.Claims().Uid)
	if err != nil {
		return ginx.Result{}, err
	}
	resp, err := h.rbacSvc.RevokeRolePermission(businessAdminCtx, &permissionv1.RevokeRolePermissionRequest{
		Id: req.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}

func (h *AccountHandler) GrantUserRole(ctx *ginx.Context, req GrantUserRoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.checkAccountPermission(businessAdminCtx, req.BizID, sess.Claims().Uid)
	if err != nil {
		return ginx.Result{}, err
	}
	resp, err := h.rbacSvc.GrantUserRole(businessAdminCtx, &permissionv1.GrantUserRoleRequest{
		UserRole: &permissionv1.UserRole{
			UserId:   req.UserID,
			RoleId:   req.Role.ID,
			RoleName: req.Role.Name,
			RoleType: domain.DefaultAccountRoleType,
		},
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: h.toUserRoleVO(resp.UserRole),
	}, nil
}

func (h *AccountHandler) RevokeUserRole(ctx *ginx.Context, req RevokeUserRoleReq, sess session.Session) (ginx.Result, error) {
	businessAdminCtx, err := h.businessAdminCtx(ctx, req.BizID)
	if err != nil {
		return ginx.Result{}, err
	}
	if err := h.checkAccountPermission(businessAdminCtx, req.BizID, sess.Claims().Uid); err != nil {
		return ginx.Result{}, err
	}
	resp, err := h.rbacSvc.RevokeRolePermission(businessAdminCtx, &permissionv1.RevokeRolePermissionRequest{
		Id: req.ID,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: resp.Success,
	}, nil
}
