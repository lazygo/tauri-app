package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/lazygo/client/framework"
)

var rnum int32 = 0

// AccessLog 访问日志记录中间件
func AccessLog(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		uri := ctx.Request().RequestURI

		errno := 0
		code := http.StatusOK
		st := time.Now()

		atomic.AddInt32(&rnum, 1)
		defer atomic.AddInt32(&rnum, -1)

		defer func() {
			rec := recover()
			if rec != nil {
				ctx.Logger().Alert("%v", rec)
				ctx.Logger().Error("%s", string(debug.Stack()))
				framework.App().ErrorHandler(fmt.Errorf("%v", rec), ctx)
				errno = 500
			}

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			ctx.Logger().Notice(
				"[pid: %d] [goroutine: %d] [sys: %.2fM] [alloc: %.2fM] [rnum: %d] [time: %.1fms] [status: %d] [errno: %d] [request_uri: %s]",
				os.Getegid(),
				runtime.NumGoroutine(),
				float64(m.Sys)/1024/1024,
				float64(m.Alloc)/1024/1024,
				atomic.LoadInt32(&rnum),
				float64(time.Now().Sub(st).Microseconds())/1000,
				code,
				errno,
				uri,
			)
		}()

		respErr := next(ctx)
		if respErr != nil {
			if he, ok := respErr.(*framework.Error); ok {
				errno = he.Errno
				code = he.Code
			}
		}

		return respErr
	}
}
