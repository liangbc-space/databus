package routers

import (
	"demo/middleware/jwt"
	"demo/routers/api"
	"demo/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	request := gin.Default()
	//request.Use(gin.Logger())
	request.NoRoute(api.Handle404)

	apiv1 := request.Group("/v1")
	apiv1.POST("/register", v1.Register)
	apiv1.POST("/login", v1.Login)

	apiv1.Use(jwt.Handle())
	apiv1.GET("/user/:id", v1.UserInfo)
	apiv1.GET("/logout", v1.Logout)

	return request
}
