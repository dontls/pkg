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

type zlog struct {
	logger  zerolog.Logger
	lwriter zerolog.LevelWriter
}

var (
	_console = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05"}
	_log     zlog
)

func init() {
	_log.lwriter = zerolog.MultiLevelWriter(_console)
	_log.logger = zerolog.New(_log.lwriter).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

func WithWriter(w *Logger, l int) {
	_log.lwriter = zerolog.MultiLevelWriter(_console, w)
	_log.logger.Output(_log.lwriter).Level(zerolog.Level(l))
}

func Writer() io.Writer {
	return _log.lwriter
}

func Log() *zerolog.Logger {
	return &_log.logger
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
