package log

import (
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger = lumberjack.Logger
type Level = zerolog.Level

var (
	logger  zerolog.Logger
	lwriter zerolog.LevelWriter
)

func init() {
	lwriter = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05"})
	logger = zerolog.New(lwriter).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

func Output(w *Logger, level int) *zerolog.Logger {
	// logger = zerolog.New(zerolog.MultiLevelWriter(_console, w)).With().Timestamp().Logger()
	lwriter = zerolog.MultiLevelWriter(lwriter, w)
	logger = logger.Output(lwriter).Level(zerolog.Level(level))
	return &logger
}

func Writer() io.Writer {
	return lwriter
}

func Log() *zerolog.Logger {
	return &logger
}

func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errstr := strings.ToLower(se.Error())
						if strings.Contains(errstr, "borken pipe") || strings.Contains(errstr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				lmsg := Log().Error().Any("error", err).Str("request", string(httpRequest))
				if brokenPipe {
					lmsg.Msg(c.Request.URL.Path)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					lmsg.Str("stack", string(debug.Stack()))
				}
				lmsg.Any("error", err).Msg(c.Request.URL.Path)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
