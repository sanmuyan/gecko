package controller

import (
	"gecko/pkg/model"
	"gecko/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sanmuyan/xpkg/xresponse"
	"github.com/sirupsen/logrus"
)

// GitlabUpdate 处理 gitlab webhook 事件
func GitlabUpdate(c *gin.Context) {
	event := &model.GitlabWebhook{}
	err := c.ShouldBindJSON(event)
	if err != nil {
		logrus.Errorf("body parse error: %v", err)
		util.Respf().Fail(xresponse.HttpBadRequest).Response(util.GinRespf(c))
		return
	}
	logrus.Infof("webhook event: %v", event.EventName)
	project, err := svc.GetProject(event.ProjectID)
	if err != nil {
		util.Respf().Fail(xresponse.HttpInternalServerError).Response(util.GinRespf(c))
		return
	}
	svc.SyncProject(project)
	util.Respf().Ok().Response(util.GinRespf(c))
}
