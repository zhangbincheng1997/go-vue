package global

import (
	"main/model"

	"github.com/gin-gonic/gin"
)

// IdentityKey ...
const IdentityKey = "username"

// GetAuthUser ...
func GetAuthUser(c *gin.Context) *model.User {
	user, _ := c.Get(IdentityKey)
	return user.(*model.User)
}
