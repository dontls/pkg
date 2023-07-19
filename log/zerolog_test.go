package log

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T) {
	Log().Debug().Uint16("port", 12).Msg("start at")
	Log().Info().Uint16("port", 63).Msg("panic")

	// WithWriter(&Logger{
	// 	Filename: "log.txt",
	// 	MaxSize:  10}, 2)
	for {
		Log().Debug().Msg(":----> debug message")
		Log().Info().Msg(":----> info message")
		Log().Error().Msg(":----> error message")
	}
}

func TestLoggerGin(t *testing.T) {
	WithWriter(&Logger{
		Filename: "log.txt",
		MaxSize:  10}, 1)
	e := gin.New()
	e.Use(Recovery(true))
	e.Use(gin.LoggerWithWriter(Writer()))
	e.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	e.GET("/over", func(ctx *gin.Context) {
		panic("over")
	})
	e.Run(":8080")
}
