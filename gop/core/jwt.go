package core

import (
	"encoding/json"
	"log"
	"main/global"
	"main/model"
	"main/model/request"
	"main/model/response"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const identityKey = global.IdentityKey

// JWT ...
func JWT() *jwt.GinJWTMiddleware {
	cfg := global.CONFIG.JWT
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         cfg.Realm,
		Key:           []byte(cfg.Key),
		Timeout:       time.Hour * time.Duration(cfg.Timeout),
		IdentityKey:   identityKey,
		TokenLookup:   cfg.TokenLookup,
		TokenHeadName: cfg.TokenHeadName,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username := claims[identityKey].(string)

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
			res := global.DB.Where("username = ? and password = ?", req.Username, req.Password).First(&user)
			if err := res.Error; err != nil {
				global.LOG.Error("登录失败", zap.Any("err", err))
				return nil, jwt.ErrFailedAuthentication
			}
			return &model.User{Username: req.Username}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			global.LOG.Info("认证用户", zap.Any("user", data.(*model.User)))
			if v, ok := data.(*model.User); ok && v.Role == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			response.Result(c, code, nil, message)
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			response.OkWithData(c, gin.H{
				"token":  token,
				"expire": t.Format(time.RFC3339),
			})
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
