package util

import (
	"github.com/gin-gonic/gin"
	"github.com/sanmuyan/xpkg/xresponse"
)

var Respf = func() *xresponse.Response {
	return xresponse.NewResponse()
}

type ginResponse struct {
	*gin.Context
}

func (r *ginResponse) SetFramework(c *xresponse.Response) {
	r.Context.JSON(200, c)
}

var GinRespf = func(c *gin.Context) *ginResponse {
	return &ginResponse{Context: c}
}
