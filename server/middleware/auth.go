package middleware

import (
	"gecko/pkg/config"
	"gecko/pkg/model"
	"gecko/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sanmuyan/dao/response"
	"github.com/sanmuyan/dao/secret"
	"time"
)

func Auth(c *gin.Context) {
	if !config.Conf.EnableAuth {
		return
	}
	var user model.OauthUser
	res := func() bool {
		reqToken := c.Request.Header.Get("token")
		if reqToken == "" {
			return false
		}
		_, err := secret.ParseToken(reqToken, config.Conf.TokenKey, &user)
		if err != nil {
			return false
		}
		if user.ExpirationTime < time.Now().Unix() {
			return false
		}
		return true
	}
	if !res() {
		util.Respf().Fail(response.HttpUnauthorized).Response(util.GinRespf(c))
		c.Abort()
		return
	}
	c.Set("userToken", user)
}
