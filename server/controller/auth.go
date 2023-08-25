package controller

import (
	"gecko/pkg/util"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	util.Respf().Ok().WithData(svc.Login()).Response(util.GinRespf(c))
}

func OauthCallback(c *gin.Context) {
	code := c.Query("code")
	util.Respf().Ok().WithData(svc.OauthCallback(code)).Response(util.GinRespf(c))
}
