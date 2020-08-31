package jwt

import (
	"demo/routers/api"
	"demo/utils"
	"github.com/gin-gonic/gin"
)

func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &api.Response{}
		jwtToken := c.GetHeader("Authorization")
		if jwtToken == "" {
			resp.Code = 400
			resp.Message = "非法操作，鉴权令牌缺失"
			resp.OutError(c)

			c.Abort()
			return
		}

		claims, err := utils.ParseToken(jwtToken)
		if err != nil {
			resp.Code = 403
			resp.Message = err.Error()
			resp.OutError(c)

			c.Abort()
			return
		}

		//	检测令牌是否被强制下线、忘记密码、修改密码等	在用户退出登录、被强制下线、忘记密码、修改密码等情况将对应的token加入到黑名单中，防止token未过期依旧可以使用
		key := "jwt_token:blacklist"
		res, err := utils.RedisClient.Exec("sismember", key, jwtToken)
		if err != nil {
			resp.Code = 403
			resp.Message = err.Error()
			resp.OutError(c)

			c.Abort()
			return
		} else if res.(int64) > 0 {
			resp.Code = 401
			resp.Message = "账号登录令牌已失效，请重新登录"
			resp.OutError(c)

			c.Abort()
			return
		}

		c.Set("loginInfo", claims)
		c.Set("jwtToken", jwtToken)

		c.Next()
	}

}
