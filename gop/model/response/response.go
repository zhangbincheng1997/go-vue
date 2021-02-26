package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response ...
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// CODE
const (
	ERROR   = 50000
	SUCCESS = 20000
)

// Result ...
func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}

// Ok ...
func Ok(c *gin.Context) {
	Result(c, SUCCESS, "ok", map[string]interface{}{})
}

// OkWithData ...
func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, "ok", data)
}

// OkWithMsg ...
func OkWithMsg(c *gin.Context, msg string) {
	Result(c, SUCCESS, msg, map[string]interface{}{})
}

// OkWithDetailed ...
func OkWithDetailed(c *gin.Context, data interface{}, msg string) {
	Result(c, SUCCESS, msg, data)
}

// Fail ...
func Fail(c *gin.Context) {
	Result(c, ERROR, "fail", map[string]interface{}{})
}

// FailWithData ...
func FailWithData(c *gin.Context, data interface{}) {
	Result(c, ERROR, "fail", data)
}

// FailWithMsg ...
func FailWithMsg(c *gin.Context, msg string) {
	Result(c, ERROR, msg, map[string]interface{}{})
}

// FailWithDetailed ...
func FailWithDetailed(c *gin.Context, data interface{}, msg string) {
	Result(c, ERROR, msg, data)
}
