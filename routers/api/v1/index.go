package v1

import (
	"demo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexGet(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"name": models.Lists(),
	})
	return
}
