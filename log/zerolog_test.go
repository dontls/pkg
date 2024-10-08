package log

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	Debug().Uint16("port", 12).Msg("start at")
	Info().Uint16("port", 63).Msg("panic")

	Output(&Logger{
		Filename: "log.txt",
		MaxSize:  10}, 2)
	for {
		Debug().Msg("message")
		Info().Msg("message")
		Error().Msg("message")
		time.Sleep(1 * time.Second)
	}
}

// func TestLoggerGin(t *testing.T) {
// 	Output(&Logger{Filename: "log.txt", MaxSize: 10}, 1)
// 	e := gin.New()
// 	gin.DefaultErrorWriter = Writer()
// 	e.Use(gin.Recovery())
// 	e.GET("/hello", func(ctx *gin.Context) {
// 		ctx.String(http.StatusOK, "hello")
// 	})
// 	e.GET("/over", func(ctx *gin.Context) {
// 		panic("over")
// 	})
// 	e.Run(":8080")
// }
