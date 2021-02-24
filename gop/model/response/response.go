package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response ...
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// CODE
const (
	ERROR   = 50000
	SUCCESS = 20000
)

// Result ...
func Result(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

// Ok ...
func Ok(c *gin.Context) {
	Result(c, SUCCESS, map[string]interface{}{}, "ok")
}

// OkWithData ...
func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "ok")
}

// OkWithMsg ...
func OkWithMsg(c *gin.Context, msg string) {
	Result(c, SUCCESS, map[string]interface{}{}, msg)
}

// OkWithDetailed ...
func OkWithDetailed(c *gin.Context, data interface{}, msg string) {
	Result(c, SUCCESS, data, msg)
}

// Fail ...
func Fail(c *gin.Context) {
	Result(c, ERROR, map[string]interface{}{}, "fail")
}

// FailWithData ...
func FailWithData(c *gin.Context, data interface{}) {
	Result(c, ERROR, data, "fail")
}

// FailWithMsg ...
func FailWithMsg(c *gin.Context, msg string) {
	Result(c, ERROR, map[string]interface{}{}, msg)
}

// FailWithDetailed ...
func FailWithDetailed(c *gin.Context, data interface{}, msg string) {
	Result(c, ERROR, data, msg)
}
