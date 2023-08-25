package controller

import (
	"gecko/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RunServer(addr string) {
	r := gin.Default()
	router(r)
	err := r.Run(addr)
	if err != nil {
		logrus.Fatal(err)
	}
}

func router(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/search", middleware.Auth, SearchCode)
		apiGroup.POST("/gitlab/update", GitlabUpdate)
		apiGroup.GET("/login", Login)
		apiGroup.GET("/oauth/callback", OauthCallback)
		apiGroup.GET("/projects", middleware.Auth, Projects)
	}
}
