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
	StatusParamErr     = 5000 // 参数错误
)

type H gin.H

type rData struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data,omitempty"`
}

// UserClaims 用户Claims
type UserClaims struct {
	UserID   uint   `json:"userId"`
	RoleID   uint   `json:"roleId"`
	UserName string `json:"username"`
}

// Context 响应
type Context struct {
	*gin.Context
	rData
}

func JSON(c *gin.Context) *Context {
	return &Context{Context: c, rData: rData{Code: StatusOK, Msg: "OK"}}
}

// JSONWriteMsg 自定义错误应答
func (c *Context) JSONWriteMsg(code int, err error) {
	c.Code = code
	c.Msg = err.Error()
	c.JSON(http.StatusOK, c.rData)
}

// WriteError 内部错误
func (c *Context) JSONWriteError(err error) {
	if err != nil {
		c.Code = StatusError
		c.Msg = err.Error()
	}
	c.JSON(http.StatusOK, c.rData)
}

// WriteData 输出json到客户端， 有data字段
func (c *Context) JSONWriteData(data any, errs ...error) {
	if len(errs) > 0 {
		c.JSONWriteError(errs[0])
		return
	}
	c.rData.Data = data
	c.JSON(http.StatusOK, c.rData)
}

// Write 输出json到客户端, 无data字段
func (c *Context) JSONWrite(h H, errs ...error) {
	if len(errs) > 0 {
		c.JSONWriteError(errs[0])
		return
	}
	h["code"] = c.Code
	h["message"] = c.Msg
	c.JSON(http.StatusOK, h)
}

// WriteData 输出json到客户端， 有data字段
func (c *Context) JSONWriteTotal(n int64, data any) {
	c.JSONWrite(H{"total": n, "data": data})
}

func (c *Context) MustBind(v any) error {
	err := c.ShouldBind(v)
	if err != nil {
		c.JSONWriteMsg(StatusParamErr, err)
	}
	return err
}

// ParamUInt uint参数
func (c *Context) MustParam(key string) string {
	idstr := c.Param(key)
	if idstr == "" {
		c.JSONWriteMsg(StatusParamErr, fmt.Errorf("%s empty", key))
	}
	return idstr
}

// ParamUInt uint参数
func (c *Context) ParamUInt(key string) uint {
	idstr := c.Param(key)
	id, _ := strconv.Atoi(idstr)
	return uint(id)
}

// ParamInt int参数
func (c *Context) ParamInt(key string) int {
	return int(c.ParamUInt(key))
}

// QueryInt int参数
func (c *Context) QueryInt(key string) int {
	idstr := c.Query(key)
	n, _ := strconv.Atoi(idstr)
	return n
}

// QueryUInt int参数
func (c *Context) QueryUInt(key string) uint {
	return uint(c.QueryInt(key))
}

// GetUser 根据Token获取用户信息
func (c *Context) GetUser() UserClaims {
	claims, _ := c.Get("claims")
	if claims == nil {
		return UserClaims{}
	}
	return claims.(UserClaims)
}
