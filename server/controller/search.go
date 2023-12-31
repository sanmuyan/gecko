package controller

import (
	"gecko/pkg/config"
	"gecko/pkg/model"
	"gecko/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sanmuyan/xpkg/xresponse"
	"github.com/sirupsen/logrus"
)

func SearchCode(c *gin.Context) {
	project := &model.Project{}
	if err := c.ShouldBindQuery(project); err != nil {
		util.Respf().Fail(xresponse.HttpBadRequest).Response(util.GinRespf(c))
		return
	}
	pageNumber, pageSize := util.GetPage(c)
	if pageNumber*pageSize > int(config.Conf.MaxSearchTotal) {
		util.Respf().Fail(xresponse.HttpBadRequest).Response(util.GinRespf(c))
		return
	}
	res, err := svc.SearchCode(project, pageNumber, pageSize)
	if err != nil {
		logrus.Errorln(err)
	}
	if res.TotalCount > config.Conf.MaxSearchTotal {
		res.TotalCount = config.Conf.MaxSearchTotal
	}
	user, err := util.GetUser(c)
	if err != nil {
		logrus.Errorln(err)
		util.Respf().Fail(xresponse.HttpInternalServerError).Response(util.GinRespf(c))
		return
	}
	if config.Conf.EnableAuth && !user.IsAdmin {
		for _, project := range res.Projects {
			if project.Visibility == "private" {
				if !svc.IsUserAccess(project, user.ID) {
					project.CodeContent = "搜索命中，但是你没有该项目权限"
				}
			}
		}
	}
	if !config.Conf.EnableCodeFullPreview {
		for _, project := range res.Projects {
			project.CodeContent = util.MinPreview(project.CodeContent)
		}
	}
	util.Respf().Ok().WithData(res).Response(util.GinRespf(c))
}

func Projects(c *gin.Context) {
	util.Respf().Ok().WithData(svc.Projects()).Response(util.GinRespf(c))
}
