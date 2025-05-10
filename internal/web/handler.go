package web

import (
	"github.com/ecodeclub/ginx"
	"github.com/ecodeclub/ginx/session"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
}

func (h *AdminHandler) RegisterRoutes(server *gin.Engine) {
	ginx.BS(h.CreateBiz)
}

func (h *AdminHandler) CreateBiz(ctx *ginx.Context, req CreateBizReq, session session.Session) (ginx.Result, error) {
	uid := session.Claims().Uid
}
