package ioc

import (
	"net/http"

	"gitee.com/flycash/permission-platform-admin/internal/web"
	"github.com/ecodeclub/ginx/session"

	"github.com/gin-gonic/gin"

	"github.com/gotomicro/ego/server/egin"
)

func initGinServer(
	sp session.Provider,
	account *web.AccountHandler,
	business *web.BusinessHandler,
	systemAdmin *web.SystemAdminHandler,
) *egin.Component {
	session.SetDefaultProvider(sp)
	res := egin.Load("server.web").Build()
	res.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world!")
	})
	res.Use(session.CheckLoginMiddleware())
	account.PrivateRoutes(res.Engine)
	business.PrivateRoutes(res.Engine)
	systemAdmin.PrivateRoutes(res.Engine)
	return res
}
