package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseWithErr(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"msg": err.Error(),
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	var respObj = gin.H{
		"msg": "success",
	}
	if data != nil {
		respObj["data"] = data
	}
	c.JSON(http.StatusOK, respObj)
}
