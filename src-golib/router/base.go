package router

import (
	"github.com/lazygo/client/app/middleware"
	"github.com/lazygo/client/framework"
)

func Init(app *framework.Application) *framework.Application {

	// 请求前去除url中 .json 后缀
	app.Pre(middleware.StripUrlSuffix)

	// 添加request_id
	app.Pre(middleware.RequestID)

	// 增加访问日志记录
	app.Use(middleware.AccessLog)

	app.Add("/", framework.NotFoundHandler)

	ApiRouter(app)

	return app
}
