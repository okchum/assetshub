package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Out(c *gin.Context, code int, msg string, data ...[]string) {
	c.JSON(http.StatusOK, gin.H{
		"message": msg,
		"code":    code,
		"data":    data,
	})
}
