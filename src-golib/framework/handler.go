package framework

import (
	"reflect"
	"strings"

	"github.com/lazygo/client/utils"
)

// Controller 转为 server.HandlerFunc
func Controller(h interface{}, methodName ...string) HandlerFunc {
	rtServ, serviceName, err := routes.Make(h)
	if err != nil {
		panic(err)
	}
	var name string
	if len(methodName) > 0 {
		name = utils.ToSnakeString(methodName[0])
	}

	return func(ctx Context) error {
		if name == "" {
			routePath := strings.TrimRight(ctx.Path(), "/")
			index := strings.LastIndex(routePath, "/")
			name = strings.TrimLeft(routePath[index:], "/")
		}

		method, ok := routes[serviceName][name]
		if !ok {
			ctx.Logger().Warn("[msg: not fount] [method name: %s]", name)
			return ErrNotFound
		}
		pReq := reflect.New(method.Request)
		req := pReq.Interface().(BaseRequest)

		defer req.Clear()

		if err = ctx.Bind(req); err != nil {
			ctx.Logger().Warn("[msg: params error] [req: %v] [err: %v]", req, err)
			return ErrBadRequest
		}
		if err = req.Verify(); err != nil {
			ctx.Logger().Warn("[msg: verify params fail] [resp: %v] [err: %v]", req, err)
			return err
		}

		pServ := reflect.New(rtServ)
		pServ.Elem().FieldByName("Ctx").Set(reflect.ValueOf(ctx))

		out := method.Method.Func.Call([]reflect.Value{pServ, pReq})
		numOut := len(out)
		if numOut == 1 {
			if ierr := out[0].Interface(); ierr != nil {
				if err = ierr.(error); err != nil {
					ctx.Logger().Warn("[msg: request fail] [req: %v] [err: %v]", req, err)
					return err
				}
			}
			return nil
		}
		if numOut == 2 {
			resp := out[0].Interface()
			if ierr := out[1].Interface(); ierr != nil {
				if err = ierr.(error); err != nil {
					ctx.Logger().Warn("[msg: request fail] [req: %v] [resp: %v] [err: %v]", req, resp, err)
					return err
				}
			}
			return ctx.Succ(resp)
		}
		ctx.Logger().Warn("[msg: method return value error] [out: %v]", out)
		return ErrInternalServerError
	}
}
