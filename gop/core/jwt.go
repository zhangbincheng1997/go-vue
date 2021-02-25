package core

import (
	"encoding/json"
	"log"
	"main/global"
	"main/model"
	"main/model/request"
	"main/model/response"
	"main/utils"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

			val, err := global.RDB.Get(username).Result()
			if err != nil {
				var user model.User
				// errors.Is(err, gorm.ErrRecordNotFound)
				if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil { // MySQL
					return nil
				}
				user.Password = "" // 敏感数据
				jsonStr, _ := json.Marshal(user)
				if err := global.RDB.Set(username, string(jsonStr), 0).Err(); err != nil { // Redis
					return nil
				}
				return &user
			}
			var user model.User
			_ = json.Unmarshal([]byte(val), &user)
			return &user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req request.LoginReq
			if err := c.ShouldBind(&req); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			var user model.User
			res := global.DB.Where("username = ? and password = ?", req.Username, utils.MD5(req.Password)).First(&user)
			if err := res.Error; err != nil {
				global.LOG.Error("登录失败", zap.Any("err", err))
				return nil, jwt.ErrFailedAuthentication
			}
			return &model.User{Username: user.Username}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
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
		log.Fatal("JWT Error:" + err.Error())
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}

//GroupAuthorizator ...
func GroupAuthorizator(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := global.GetAuthUser(c)
		global.LOG.Info("认证用户", zap.Any("user", user))
		if role == user.Role {
			c.Next()
			return
		}
		c.Abort()
		response.FailWithMsg(c, jwt.ErrForbidden.Error())
	}
}
