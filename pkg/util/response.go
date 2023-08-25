package util

import (
	"github.com/gin-gonic/gin"
	"github.com/sanmuyan/dao/response"
)

var Respf = func() *response.Response {
	return response.NewResponse()
}

type ginResponse struct {
	*gin.Context
}

func (r *ginResponse) SetFramework(c *response.Response) {
	r.Context.JSON(200, c)
}

var GinRespf = func(c *gin.Context) *ginResponse {
	return &ginResponse{Context: c}
}
