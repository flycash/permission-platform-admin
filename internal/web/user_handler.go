package web

import (
	"fmt"
	"github.com/ecodeclub/ginx/session"
	"time"

	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ginx"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
	"github.com/redis/go-redis/v9"
)

const (
	defaultTimeout = 30 * 24 * time.Hour
	permissionName = "permission"
)

type UserHandler struct {
	*BaseHandler
	client redis.Cmdable
	logger *elog.Component
	sp     session.Provider
}

func NewUserHandler(handler *BaseHandler, client redis.Cmdable) *UserHandler {
	return &UserHandler{BaseHandler: handler, logger: elog.DefaultLogger, client: client}
}

func (h *UserHandler) PublicRoutes(server *gin.Engine) {

}

func (h *UserHandler) LoginDemo(ctx *ginx.Context, req LoginReq) (ginx.Result, error) {
	// todo 通过账号密码找到用户id，模拟就直接使用传进来的id作为uid了
	uid := req.ID
	bizid := req.BizID
	resp, err := h.rbacSvc.GetAllPermissions(ctx, &permissionv1.GetAllPermissionsRequest{
		BizId:  bizid,
		UserId: uid,
	})
	if err != nil {
		return ginx.Result{}, err
	}
	permissionList := slice.Map(resp.GetUserPermissions(), func(idx int, src *permissionv1.UserPermission) domain.UserPermission {
		return domain.UserPermission{
			BizID:  bizid,
			UserID: uid,
			Permission: domain.Permission{
				Resource: domain.Resource{
					Type: src.ResourceType,
					Key:  src.ResourceKey,
				},
				Action: src.PermissionAction,
			},
			Effect: domain.Effect(src.Effect),
		}
	})
	res, err := session.NewSessionBuilder(ctx, uid).SetSessData(map[string]any{
		permissionName: permissionList,
	}).Build()
	if err != nil {
		return ginx.Result{}, fmt.Errorf("设置session失败：%w", err)
	}
	_, err = h.client.Set(ctx, fmt.Sprintf("%d", uid), res.Claims().SSID, defaultTimeout).Result()
	if err != nil {
		return ginx.Result{}, fmt.Errorf("设置用户权限缓存和token失败 %w", err)
	}
	return ginx.Result{}, nil
}
