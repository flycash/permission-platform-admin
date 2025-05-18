//go:build wireinject

package ioc

import (
	"gitee.com/flycash/permission-platform-admin/internal/web"
	"github.com/google/wire"
)

func InitApp() (*App, error) {
	wire.Build(

		InitRedis,
		InitSession,

		InitRBACClient,
		InitPermissionClient,
		InitBaseHandler,

		// 账号管理API
		web.NewAccountHandler,

		// 业务方使用的管理后台API
		web.NewBusinessHandler,

		// 权限平台系统管理员使用的管理后台API
		web.NewSystemAdminHandler,

		initGinServer,

		wire.Struct(new(App), "*"),
	)
	return new(App), nil
}
