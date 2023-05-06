package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	c *gin.Context
}

func newResponse(c *gin.Context) *response {
	return &response{
		c: c,
	}
}

func (r *response) success() {
	r.c.JSON(
		http.StatusOK,
		gin.H{
			"msg": "ok",
		},
	)
}

func (r *response) successWithAny(data any) {
	r.c.JSON(
		http.StatusOK,
		gin.H{
			"msg":  "ok",
			"data": data,
		},
	)
}

func (r *response) fail(msg string) {
	r.c.JSON(
		http.StatusOK,
		gin.H{
			"msg": msg,
		},
	)
}
