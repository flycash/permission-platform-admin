package web

import (
	"errors"
	"fmt"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/ecodeclub/ginx"
	"github.com/ecodeclub/ginx/session"
	"google.golang.org/grpc/metadata"
)

type BizHandler struct {
	rbacSvc permissionv1.RBACServiceClient
	svc     permissionv1.PermissionServiceClient
	// biz_id = 1 的 token
	adminToken string
}

// CreateRole 创建业务内的角色，比如说店长
func (b *BizHandler) CreateRole(ctx *ginx.Context, req CreateRoleRequest, sess session.Session) (ginx.Result, error) {
	gctx := metadata.AppendToOutgoingContext(ctx, "Authorization", b.adminToken)
	bizResp, err := b.rbacSvc.GetBusinessConfig(gctx, &permissionv1.GetBusinessConfigRequest{
		Id: req.BizID,
	})
	if err != nil {
		return ginx.Result{}, err
	}

	gctx = metadata.AppendToOutgoingContext(ctx, "Authorization", bizResp.Config.Token)
	// biz_id = 2
	uid := sess.Claims().Uid
	resp, err := b.svc.CheckPermission(gctx, &permissionv1.CheckPermissionRequest{
		Uid: uid,
		Permission: &permissionv1.Permission{
			ResourceKey: fmt.Sprintf("/admin/role/%d", req.BizID),
		},
	})
	if err != nil {
		return ginx.Result{}, err
	}
	if !resp.Allowed {
		return ginx.Result{}, errors.New("没有权限")
	}

	roleResp, err := b.rbacSvc.CreateRole(gctx, &permissionv1.CreateRoleRequest{
		Role: &permissionv1.Role{
			Name: req.Name,
		},
	})
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Data: roleResp.Role.Id,
	}, nil
}

type CreateRoleRequest struct {
	BizID int64
	Name  string
}
