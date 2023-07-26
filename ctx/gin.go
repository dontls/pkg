package ctx

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	StatusFail         = -1  // 失败
	StatusOK           = 0   // 成功
	StatusError        = 500 // 错误
	StatusLoginExpired = 401 // 登录过期
	StatusForbidden    = 403 // 无权限
)

type rdata struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

// Context 响应
type Context struct {
	rdata
}

// JSON 响应
func JSON(status int) *Context {
	c := &Context{}
	c.Code = status
	switch status {
	case StatusOK:
		c.Msg = "success"
	case StatusFail:
		c.Msg = "failed"
	case StatusForbidden:
		c.Msg = "forbidden"
	}
	return c
}

// SetMsg 设置消息体的内容int
func (o *Context) SetMsg(msg string) *Context {
	o.Msg = msg
	return o
}

// SetCode 设置消息体的编码
func (o *Context) SetCode(code int) *Context {
	o.Code = code
	return o
}

// WriteData 输出json到客户端， 有data字段
func (o *Context) WriteData(data interface{}, c *gin.Context) {
	o.Data = data
	c.JSON(http.StatusOK, o.rdata)
}

// Write 输出json到客户端, 无data字段
func (o *Context) Write(h gin.H, c *gin.Context) {
	h["code"] = o.Code
	h["msg"] = o.Msg
	c.JSON(http.StatusOK, h)
}

// JSONWrite
func JSONWrite(h gin.H, c *gin.Context) {
	JSON(StatusOK).Write(h, c)
}

// JSONWriteTotal
func JSONWriteTotal(total int64, data interface{}, c *gin.Context) {
	JSON(StatusOK).Write(gin.H{"total": total, "data": data}, c)
}

// JSONWriteData
func JSONWriteData(data interface{}, c *gin.Context) {
	JSON(StatusOK).WriteData(data, c)
}

// JSONWriteError 错误应答
func JSONWriteError(err error, c *gin.Context) {
	ctx := JSON(StatusError)
	if err != nil {
		ctx.SetMsg(err.Error())
	}
	ctx.WriteData(nil, c)
}

// ParamUInt uint参数
func ParamUInt(c *gin.Context, key string) uint {
	idstr := c.Param(key)
	id, _ := strconv.Atoi(idstr)
	return uint(id)
}

// ParamInt int参数
func ParamInt(c *gin.Context, key string) int {
	return int(ParamUInt(c, key))
}

// QueryInt int参数
func QueryInt(c *gin.Context, key string) (int, error) {
	idstr := c.Query(key)
	return strconv.Atoi(idstr)
}

// QueryUInt int参数
func QueryUInt(c *gin.Context, key string) (uint, error) {
	idstr := c.Query(key)
	id, err := strconv.Atoi(idstr)
	return uint(id), err
}

// QueryUInt64 int参数
func QueryUInt64(c *gin.Context, key string) (uint64, error) {
	idstr := c.Query(key)
	id, err := strconv.Atoi(idstr)
	return uint64(id), err
}
