package controller

import (
	request "github.com/lazygo/client/app/request/system"
	"github.com/lazygo/client/framework"
	"github.com/lazygo/client/pkg/httpd"
)

type SystemController struct {
	Ctx framework.Context
}

func (c *SystemController) Version(req *request.SystemVersionRequest) (*request.SystemVersionResponse, error) {
	resp := &request.SystemVersionResponse{
		Version: framework.App().Build,
	}

	app := framework.App()
	resp.Build = app.Build
	resp.Debug = app.Debug

	return resp, nil
}

func (s *SystemController) Logger(req *request.SystemLoggerRequest) (*request.SystemLoggerResponse, error) {
	resp := &request.SystemLoggerResponse{}

	switch req.Level {
	case "error":
		s.Ctx.Logger().Error("[msg: client] %s", req.Content)
	case "warning":
		s.Ctx.Logger().Warn("[msg: client] %s", req.Content)
	case "debug":
		fallthrough
	default:
		s.Ctx.Logger().Debug("[msg: client] %s", req.Content)
	}
	return resp, nil
}

func (s *SystemController) Httpd(req *request.SystemHttpdRequest) (*request.SystemHttpdResponse, error) {
	resp := &request.SystemHttpdResponse{}
	resp.BeforeState = httpd.Starting()
	if req.Method == "start" {
		httpd.AsyncStart(req.Addr)
	}
	if req.Method == "stop" {
		httpd.Stop()
	}
	resp.AfterState = httpd.Starting()
	return resp, nil
}
