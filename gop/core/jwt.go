package core

import (
	"encoding/json"
	"errors"
	"main/global"
	"main/model"
	"main/model/request"
	"main/model/response"
	"main/utils"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JWT ...
func JWT() *jwt.GinJWTMiddleware {
	cfg := global.CONFIG.JWT
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:     []byte(cfg.Key),
		Timeout: time.Hour * time.Duration(cfg.Timeout),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username := claims[jwt.IdentityKey].(string)
			return username
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req request.LoginReq
			if err := c.ShouldBind(&req); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			var user model.User
			err := global.DB.Where("username = ? and password = ?", req.Username, utils.MD5(req.Password)).First(&user).Error
			if err != nil {
				global.LOG.Errorf("登录失败：%v", err)
				return nil, jwt.ErrFailedAuthentication
			}
			return &model.User{Username: user.Username}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if username, ok := data.(string); ok && username != "" {
				return true
			}
			return true // GroupAuthorizator
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			response.FailWithMsg(c, message)
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			response.OkWithData(c, gin.H{
				"token":  token,
				"expire": t.Format(time.RFC3339),
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			response.Ok(c)
		},
	})
	if err != nil {
		global.LOG.Errorf("JWT: %v", err)
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		global.LOG.Errorf("MiddlewareInit: %v", errInit)
	}
	return authMiddleware
}

//GroupAuthorizator ...
func GroupAuthorizator(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := global.GetAuthUser(c)

		var user model.User
		val, err := global.RDB.Get(username).Result()
		if err != nil {
			if errors.Is(global.DB.Where("username = ?", username).First(&user).Error, gorm.ErrRecordNotFound) { // MySQL
				c.Abort()
				response.FailWithMsg(c, jwt.ErrForbidden.Error())
				return
			}
			user.Password = "" // 敏感数据
			jsonStr, _ := json.Marshal(user)
			if err := global.RDB.Set(username, string(jsonStr), 0).Err(); err != nil { // Redis
				c.Abort()
				response.FailWithMsg(c, jwt.ErrForbidden.Error())
				return
			}
		} else {
			_ = json.Unmarshal([]byte(val), &user)
		}
		global.LOG.Infof("认证用户：%v", user)
		if role == user.Role {
			c.Next()
			return
		}
		c.Abort()
		response.FailWithMsg(c, jwt.ErrForbidden.Error())
	}
}
