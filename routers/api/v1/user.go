package v1

import (
	"databus/models"
	"databus/routers/api"
	"databus/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func UserInfo(c *gin.Context) {
	loginInfo := c.MustGet("loginInfo").(*utils.JwtClaims)
	resp := new(api.Response)
	resp.Data = loginInfo
	resp.OutSuccess(c)
	return
}

func Register(c *gin.Context) {
	resp := &api.Response{}

	username := c.PostForm("name")
	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	usernameLen := strings.Count(username, "") - 1
	if usernameLen < 3 || usernameLen > 16 {
		resp.Message = "注册用户名必须是3-16不为空的字符串"
		resp.OutError(c)
		return
	}

	if mobile == "" || !utils.ValidateMobile(mobile) {
		resp.Message = "注册手机号不正确"
		resp.OutError(c)
		return
	}

	passwordLen := strings.Count(password, "") - 1
	if passwordLen < 6 {
		resp.Message = "注册密码长度不能低于6位"
		resp.OutError(c)
		return
	}

	user := &models.User{
		Username: username,
		Mobile:   mobile,
		Password: utils.Md5(password),
	}
	err := user.Insert()
	if err != nil {
		resp.Message = fmt.Sprintf("%s", err)
		resp.OutError(c)
		return
	}

	resp.OutSuccess(c)
	return
}

func Login(c *gin.Context) {
	resp := &api.Response{}

	mobile := c.PostForm("mobile")
	password := c.PostForm("password")

	if mobile == "" || !utils.ValidateMobile(mobile) {
		resp.Message = "手机号不正确"
		resp.OutError(c)
		return
	}

	if password == "" {
		resp.Message = "密码不能为空哟"
		resp.OutError(c)
		return
	}

	token, err := models.Login(mobile, password)
	if err != nil {
		resp.Message = err.Error()
		resp.OutError(c)
		return
	}

	resp.Data = map[string]string{"token": token}
	resp.OutSuccess(c)
	return
}

func Logout(c *gin.Context) {
	jwtToken := c.MustGet("jwtToken").(string)

	resp := new(api.Response)
	err := models.Logout(jwtToken)
	if err != nil {
		resp.Message = err.Error()
		resp.OutError(c)
		return
	}

	resp.OutSuccess(c)
	return
}
