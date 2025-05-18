package ioc

import (
	"gitee.com/flycash/permission-platform-admin/internal/web"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/econf"
)

func InitRBACClient() permissionv1.RBACServiceClient {
	return permissionv1.NewRBACServiceClient(egrpc.Load("server.grpc.rbac").Build())
}

func InitPermissionClient() permissionv1.PermissionServiceClient {
	return permissionv1.NewPermissionServiceClient(egrpc.Load("server.grpc.rbac").Build())
}

func InitBaseHandler(
	rbacSvc permissionv1.RBACServiceClient,
	permSvc permissionv1.PermissionServiceClient,
) *web.BaseHandler {
	return web.NewBaseHandler(rbacSvc, permSvc, econf.GetString("adminToken"))
}
