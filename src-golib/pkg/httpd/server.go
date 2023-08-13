package httpd

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/lazygo/client/framework"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		req, err := framework.NewRequestWithContext(ctx, r.RequestURI, r.Body)
		if err != nil {
			// bad request
		}

		framework.App().Exec(req, framework.NewResponseWriter(w))
	})
}

type LocalServer struct {
	Server   *http.Server
	Starting atomicBool
}

var localServer LocalServer

func AsyncStart(addr string) {
	if localServer.Starting.isSet() {
		return
	}
	localServer.Starting.setTrue()
	go func() {
		defer func() {
			localServer.Starting.setFalse()
		}()
		localServer.Server = &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		}
		if err := localServer.Server.ListenAndServe(); err != http.ErrServerClosed {
			framework.GlobalLog.Warn("[msg: listen local server fail] [err: %v]", err)
		}
	}()
}

func Stop() {
	localServer.Starting.setFalse()
	go func() {
		err := localServer.Server.Close()
		if err != nil {
			framework.GlobalLog.Warn("[msg: stop local server fail] [err: %v]", err)
		}
	}()
}

func Starting() int {
	return int(localServer.Starting)
}

type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }
func (b *atomicBool) setTrue()    { atomic.StoreInt32((*int32)(b), 1) }
func (b *atomicBool) setFalse()   { atomic.StoreInt32((*int32)(b), 0) }
