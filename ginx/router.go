package ginx

import (
	"github.com/gin-gonic/gin"
)

// HandlerFunc 自定义处理函数
type HandlerFunc func(*Context)

type IRouter interface {
	Routers(*RouterGroup)
}

// RouterGroup 扩展自 gin.RouterGroup
type RouterGroup struct {
	*gin.RouterGroup
}

// NewGroup 创建路由组
func (rg *RouterGroup) Group(relativePath string) *RouterGroup {
	return &RouterGroup{
		RouterGroup: rg.RouterGroup.Group(relativePath),
	}
}

// Register 注册子路由组
func (rg *RouterGroup) Routers(relativePath string, r IRouter) {
	r.Routers(rg.Group(relativePath))
}

// Handle 统一处理方法
func (rg *RouterGroup) Handle(method, path string, handlers HandlerFunc) {
	// 调用Gin原始方法
	rg.RouterGroup.Handle(method, path, func(ctx *gin.Context) {
		handlers(JSON(ctx))
	})
}

// HTTP方法封装 =========================================================

// POST 添加POST路由
func (rg *RouterGroup) POST(relativePath string, handlers HandlerFunc) {
	rg.Handle("POST", relativePath, handlers)
}

// GET 添加GET路由
func (rg *RouterGroup) GET(relativePath string, handlers HandlerFunc) {
	rg.Handle("GET", relativePath, handlers)
}

// PUT 添加PUT路由
func (rg *RouterGroup) PUT(relativePath string, handlers HandlerFunc) {
	rg.Handle("PUT", relativePath, handlers)
}

// DELETE 添加DELETE路由
func (rg *RouterGroup) DELETE(relativePath string, handlers HandlerFunc) {
	rg.Handle("DELETE", relativePath, handlers)
}

// PATCH 添加PATCH路由
func (rg *RouterGroup) PATCH(relativePath string, handlers HandlerFunc) {
	rg.Handle("PATCH", relativePath, handlers)
}

// OPTIONS 添加OPTIONS路由
func (rg *RouterGroup) OPTIONS(relativePath string, handlers HandlerFunc) {
	rg.Handle("OPTIONS", relativePath, handlers)
}

// HEAD 添加HEAD路由
func (rg *RouterGroup) HEAD(relativePath string, handlers HandlerFunc) {
	rg.Handle("HEAD", relativePath, handlers)
}

// ANY 添加处理所有HTTP方法的路由
func (rg *RouterGroup) ANY(relativePath string, handlers HandlerFunc) {
	rg.Handle("", relativePath, handlers)
}
