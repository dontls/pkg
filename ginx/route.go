package ginx

import "github.com/gin-gonic/gin"

type RouterGroup struct {
	*gin.RouterGroup
}

type HandlerFunc func(*Context)

func (rg *RouterGroup) Group(relativePath string) *RouterGroup {
	return &RouterGroup{RouterGroup: rg.RouterGroup.Group(relativePath)}
}

func (rg *RouterGroup) Register(relativePath string, r RouterFunc) {
	r(rg.Group(relativePath))
}

func (rg *RouterGroup) hookHandler(ctx *gin.Context, handler HandlerFunc) {
	handler(&Context{Context: ctx})
}

func (rg *RouterGroup) POST(relativePath string, handler HandlerFunc) {
	rg.RouterGroup.POST(relativePath, func(ctx *gin.Context) {
		rg.hookHandler(ctx, handler)
	})
}

func (rg *RouterGroup) GET(relativePath string, handler HandlerFunc) {
	rg.RouterGroup.GET(relativePath, func(ctx *gin.Context) {
		rg.hookHandler(ctx, handler)
	})
}

func (rg *RouterGroup) PUT(relativePath string, handler HandlerFunc) {
	rg.RouterGroup.PUT(relativePath, func(ctx *gin.Context) {
		rg.hookHandler(ctx, handler)
	})
}

func (rg *RouterGroup) DELETE(relativePath string, handler HandlerFunc) {
	rg.RouterGroup.DELETE(relativePath, func(ctx *gin.Context) {
		rg.hookHandler(ctx, handler)
	})
}
