package router

import (
	"github.com/lazygo/client/app/controller"
	"github.com/lazygo/client/framework"
)

func ApiRouter(app *framework.Application) {
	g := app.Group("system")
	{
		g.Add("version", framework.Controller(controller.SystemController{}))
		g.Add("logger", framework.Controller(controller.SystemController{}))
		g.Add("httpd", framework.Controller(controller.SystemController{}))
	}
}
