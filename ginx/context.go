package ginx

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	StatusOK           = 200  // 成功
	StatusLoginExpired = 401  // 登录过期
	StatusForbidden    = 403  // 无权限
	StatusError        = 500  // 错误
	StatusParamErr     = 4000 // 参数错误
	StatusDBErr        = 4001 // 数据操作
)

type rData struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

// Context 响应
type Context struct {
	*gin.Context
	rData
}

func JSON(c *gin.Context) *Context {
	return &Context{Context: c, rData: rData{Code: StatusOK, Msg: "OK"}}
}

// SetMsg 设置消息体的内容int
func (c *Context) SetMsg(msg string) *Context {
	c.Msg = msg
	return c
}

// SetCode 设置消息体的编码
func (c *Context) SetCode(code int) *Context {
	c.Code = code
	return c
}

// Write 输出json到客户端, 无data字段
func (c *Context) Write(h gin.H) {
	h["code"] = c.Code
	h["message"] = c.Msg
	c.JSON(http.StatusOK, h)
}

// WriteData 输出json到客户端， 有data字段
func (c *Context) WriteData(data interface{}, errs ...error) {
	if len(errs) > 0 && errs[0] != nil {
		c.Code = StatusError
		c.Msg = errs[0].Error()
	} else {
		c.rData.Data = data
	}
	c.JSON(http.StatusOK, c.rData)
}

// WriteError db 错误应答
func (c *Context) WriteError(err error) {
	if err != nil {
		c.Code = StatusDBErr
		c.Msg = err.Error()
	}
	c.JSON(http.StatusOK, c.rData)
}

func (c *Context) WriteParamError(err error) {
	c.Code = StatusParamErr
	c.Msg = err.Error()
	c.JSON(http.StatusOK, c.rData)
}

// WriteData 输出json到客户端， 有data字段
func (c *Context) WriteTotal(n int64, data interface{}) {
	c.Write(gin.H{"total": n, "data": data})
}

func MustBind(c *gin.Context, v interface{}) (*Context, error) {
	ctx := JSON(c)
	err := ctx.ShouldBind(v)
	if err != nil {
		ctx.WriteParamError(err)
	}
	return ctx, err
}

// ParamUInt uint参数
func MustParam(c *gin.Context, key string) (*Context, string) {
	ctx := JSON(c)
	idstr := ctx.Param(key)
	if idstr == "" {
		ctx.WriteParamError(fmt.Errorf("%s empty", key))
	}
	return ctx, idstr
}

// ParamUInt uint参数
func (c *Context) ParamUInt(key string) (*Context, uint) {
	idstr := c.Param(key)
	id, _ := strconv.Atoi(idstr)
	return c, uint(id)
}

// ParamInt int参数
func (c *Context) ParamInt(key string) (*Context, int) {
	_, n := c.ParamUInt(key)
	return c, int(n)
}

// QueryInt int参数
func (c *Context) QueryInt(key string) (*Context, int) {
	idstr := c.Query(key)
	n, _ := strconv.Atoi(idstr)
	return c, n
}

// QueryUInt int参数
func (c *Context) QueryUInt(key string) (*Context, uint) {
	_, n := c.QueryInt(key)
	return c, uint(n)
}
