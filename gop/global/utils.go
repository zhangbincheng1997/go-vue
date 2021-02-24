package global

import (
	"main/model"

	"github.com/gin-gonic/gin"
)

// IdentityKey ...
const IdentityKey = "username"

// GetAuthUser ...
func GetAuthUser(c *gin.Context) *model.User {
	claims, exists := c.Get(IdentityKey)
	if !exists {
		LOG.Error("JWT解析错误！！！")
		return nil
	}
	return claims.(*model.User)
}
