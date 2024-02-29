package ctx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HandleFunc func(*gin.RouterGroup)

type router struct {
	key     string
	handler HandleFunc
}

var (
	routers    []router
	jwtRouters []router
	apiRouters []router
)

func Register(root string, h HandleFunc) {
	routers = append(routers, router{key: root, handler: h})
}

func RegisterJWT(root string, h HandleFunc) {
	jwtRouters = append(jwtRouters, router{key: root, handler: h})
}

func RegisterAPI(root string, h HandleFunc) {
	apiRouters = append(apiRouters, router{key: root, handler: h})
}

// 顺序，1->普通url， 2->jwtUrl, 3->tokenUrl
func Use(r ...*gin.RouterGroup) {
	if len(r) > 0 {
		for _, v := range routers {
			v.handler(r[0].Group(v.key))
		}
	}
	if len(r) > 1 {
		for _, v := range jwtRouters {
			v.handler(r[1].Group(v.key))
		}
	}
	if len(r) > 2 {
		for _, v := range apiRouters {
			v.handler(r[2].Group(v.key))
		}
	}
}

var s *http.Server

func ListenAndServe(e *gin.Engine, port, timeout int) *http.Server {
	s = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        e,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()
	return s
}

func Release() error {
	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx)
		// catching ctx.Done(). timeout of 5 seconds.
		<-ctx.Done()
	}
	return nil
}
