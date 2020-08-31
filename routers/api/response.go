package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Response struct {
	Code    int
	Data    interface{}
	Message string
}

func (resp *Response) OutSuccess(c *gin.Context) {

	if resp.Code == 0 {
		resp.Code = 200
	}

	if resp.Message == "" {
		resp.Message = "success"
	}

	h := gin.H{
		"code":      resp.Code,
		"data":      resp.Data,
		"message":   resp.Message,
		"timestamp": time.Now().UnixNano() / 1e6,
	}

	c.JSON(http.StatusOK, h)
}

func (resp *Response) OutError(c *gin.Context) {

	if resp.Code == 0 {
		resp.Code = 300
	}

	if resp.Message == "" {
		resp.Message = "failed"
	}

	h := gin.H{
		"code":      resp.Code,
		"data":      resp.Data,
		"message":   resp.Message,
		"timestamp": time.Now().UnixNano() / 1e6,
	}

	c.JSON(http.StatusOK, h)
}

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func HandleMessage(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": message,
	})
}
