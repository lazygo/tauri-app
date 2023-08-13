package main

/*
typedef int Callback(char*, char*);
static inline void CallFunc(Callback* fn, char* uri, char* data) {
	fn(uri, data);
}
*/
import "C"
import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/lazygo/client/config"
	"github.com/lazygo/client/framework"
	"github.com/lazygo/client/router"
)

func init() {

	err := config.Init()
	if err != nil {
		log.Println(string(debug.Stack()))
		log.Panicf("[msg: init config fail] [err: %v]", err)
	}
	app := framework.App()
	app.Debug = config.AppConfig.Debug

	router.Init(app).Initialized = true
}

func main() {
	//cb := func(uri string, data string) {}
	//
	//httpd.InitDefaultHttpMux()
	//// 初始化app
	//err := framework.InitApp(cb)
	//fmt.Println(err)
	//
	resp := Request(C.CString("/system/version"), C.CString("{}"))
	fmt.Println(C.GoString(resp))
	//time.Sleep(time.Hour)
}

//export InitApp
/*func InitApp(callback *C.Callback) *C.char {
	data := framework.Map{
		"code":       200,
		"errno":      0,
		"message":    "init app success",
		"data":       "",
		"request_id": 0,
	}

	cb := func(uri string, data string) {
		C.CallFunc(callback, C.CString(uri), C.CString(data))
	}

	httpd.InitDefaultHttpMux()
	// 初始化app
	err := framework.InitApp(cb)
	if err != nil {
		data["code"] = "0"
		data["errno"] = "1"
		data["message"] = "init app fail"
		data["data"] = err
	}

	resp, err := json.Marshal(data)
	if err != nil {
		data["code"] = "0"
		data["errno"] = "1"
		data["message"] = "init app fail"
		data["data"] = err
	}

	return C.CString(string(resp))
}*/

//export Request
func Request(uri *C.char, data *C.char) *C.char {
	ctx := context.Background()
	r, err := framework.NewRequestWithContext(ctx, C.GoString(uri), strings.NewReader(C.GoString(data)))
	if err != nil {
		// bad request
	}

	w := &bytes.Buffer{}
	framework.App().Exec(r, framework.NewResponseWriter(w))

	return C.CString(w.String())
}
