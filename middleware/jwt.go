package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main.go/global"
	"main.go/model/common/response"
	"main.go/service"
)

var manageAdminUserTokenService = service.ServiceGroupApp.ManageServiceGroup.ManageAdminUserTokenService
var mallUserTokenService = service.ServiceGroupApp.MallServiceGroup.MallUserTokenService

func AdminJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.FailWithDetailed(nil, "未登录或登陆失效", c)
			c.Abort()
			return
		}
		tokenClaims, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(global.GVA_CONFIG.Redis.SecretKey), nil
		})
		// err, mallAdminUserToken := manageAdminUserTokenService.ExistAdminToken(token)
		if err != nil {
			response.FailWithDetailed(nil, "未登录或登陆失效", c)
			c.Abort()
			return
		}
		if !tokenClaims.Valid {
			response.FailWithDetailed(nil, "授权已过期", c)
			c.Abort()
			return
		}
		if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				// response.FailWithCode(response.ResponseCodeUnauthorized, ctx)
				// ctx.AbortWithStatus(http.StatusUnauthorized)
				response.FailWithDetailed(nil, "授权已过期", c)
				c.Abort()
				return
			}

			userName, ok := claims["userName"].(string)
			if ok {
				// 通过key获取redis
				key := userName
				redisToken := global.GVA_REDIS.Get(key).Val()
				if token != redisToken {
					// response.FailWithCode(response.ResponseCodeUnauthorized, ctx)
					response.FailWithDetailed(nil, "未登录或登陆失效", c)
					c.Abort()
					return
				}
				c.Set("userName", userName)
			} else {
				response.FailWithDetailed(nil, "未登录或登陆失效", c)
				c.Abort()
				return
			}
		}

		c.Next()
	}

}

func UserJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.UnLogin(nil, c)
			c.Abort()
			return
		}
		err, mallUserToken := mallUserTokenService.ExistUserToken(token)
		if err != nil {
			response.UnLogin(nil, c)
			c.Abort()
			return
		}
		if time.Now().After(mallUserToken.ExpireTime) {
			response.FailWithDetailed(nil, "授权已过期", c)
			err = mallUserTokenService.DeleteMallUserToken(token)
			if err != nil {
				response.FailWithDetailed(nil, "未登录或登陆失效", c)
				return
			}
			c.Abort()
			return
		}
		c.Next()
	}

}
