package gctx

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

type Context = gin.Context
type H = gin.H

// Response Contxt 响应
type Ctx struct {
	*Context
	rData
}

func NewCtx(c *Context) *Ctx {
	return &Ctx{Context: c, rData: rData{Code: StatusOK, Msg: "OK"}}
}

// WriteData 输出json到客户端， 有data字段
func JSONWriteData(c *Context, data any, errs ...error) {
	NewCtx(c).JSONWriteData(data, errs...)
}

// WriteError 内部错误
func JSONWrite(c *Context, h H, errs ...error) {
	NewCtx(c).JSONWrite(h, errs...)
}

// WriteError 内部错误
func JSONWriteError(c *Context, err error) {
	NewCtx(c).JSONWriteError(err)
}

// JSONWriteMsg 自定义错误应答
func (c *Ctx) JSONWriteMsg(code int, err error) {
	c.Code = code
	c.Msg = err.Error()
	c.JSON(http.StatusOK, c.rData)
}

// WriteError 内部错误
func (c *Ctx) JSONWriteError(err error) {
	if err != nil {
		c.Code = StatusError
		c.Msg = err.Error()
	}
	c.JSON(http.StatusOK, c.rData)
}

// WriteData 输出json到客户端， 有data字段
func (c *Ctx) JSONWriteData(data any, errs ...error) {
	if len(errs) > 0 && errs[0] != nil {
		c.JSONWriteError(errs[0])
		return
	}
	c.rData.Data = data
	c.JSON(http.StatusOK, c.rData)
}

// Write 输出json到客户端, 无data字段
func (c *Ctx) JSONWrite(h gin.H, errs ...error) {
	if len(errs) > 0 && errs[0] != nil {
		c.JSONWriteError(errs[0])
		return
	}
	h["code"] = c.Code
	h["message"] = c.Msg
	c.JSON(http.StatusOK, h)
}

// WriteData 输出json到客户端， 有data字段
func (c *Ctx) JSONWriteTotal(n int64, data any) {
	c.JSONWrite(gin.H{"total": n, "data": data})
}

func Bind(c *Context, v any) *Ctx {
	ctx := &Ctx{Context: c}
	ctx.ShouldBind(v)
	return ctx
}

func MustBind(c *Context, v any) (*Ctx, error) {
	ctx := &Ctx{Context: c}
	err := ctx.ShouldBind(v)
	if err != nil {
		ctx.JSONWriteMsg(StatusParamErr, err)
	}
	return ctx, err
}

// ParamUInt uint参数
func MustParam(c *Context, key string) (*Ctx, string) {
	ctx := &Ctx{Context: c}
	idstr := c.Param(key)
	if idstr == "" {
		ctx.JSONWriteMsg(StatusParamErr, fmt.Errorf("%s empty", key))
	}
	return ctx, idstr
}

// ParamUInt uint参数
func ParamUInt(c *Context, key string) (*Ctx, uint) {
	ctx := &Ctx{Context: c}
	idstr := ctx.Param(key)
	id, _ := strconv.Atoi(idstr)
	return ctx, uint(id)
}

// ParamInt int参数
func ParamInt(c *Context, key string) (*Ctx, int) {
	ctx, v := ParamUInt(c, key)
	return ctx, int(v)
}

// QueryInt int参数
func QueryInt(c *Context, key string) (*Ctx, int) {
	ctx := &Ctx{Context: c}
	idstr := ctx.Query(key)
	n, _ := strconv.Atoi(idstr)
	return ctx, n
}

// QueryUInt int参数
func QueryUInt(c *Context, key string) (*Ctx, uint) {
	ctx, v := QueryInt(c, key)
	return ctx, uint(v)
}

// GetUser 根据Token获取用户信息
func (c *Ctx) GetUser() UserClaims {
	claims, _ := c.Get("claims")
	if claims == nil {
		return UserClaims{}
	}
	return claims.(UserClaims)
}
