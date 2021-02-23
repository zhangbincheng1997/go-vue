package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response ...
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// CODE
const (
	ERROR   = 50000
	SUCCESS = 20000
)

// Result ...
func Result(c *gin.Context, code int, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		message,
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

// OkWithMessage ...
func OkWithMessage(c *gin.Context, message string) {
	Result(c, SUCCESS, map[string]interface{}{}, message)
}

// OkWithDetailed ...
func OkWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, SUCCESS, data, message)
}

// Fail ...
func Fail(c *gin.Context) {
	Result(c, ERROR, map[string]interface{}{}, "fail")
}

// FailWithData ...
func FailWithData(c *gin.Context, data interface{}) {
	Result(c, ERROR, data, "fail")
}

// FailWithMessage ...
func FailWithMessage(c *gin.Context, message string) {
	Result(c, ERROR, map[string]interface{}{}, message)
}

// FailWithDetailed ...
func FailWithDetailed(c *gin.Context, code int, data interface{}, message string) {
	Result(c, code, data, message)
}
