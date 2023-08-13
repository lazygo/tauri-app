package framework

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lazygo/lazygo/server"
)

type Map map[string]interface{}
type ErrorHandler func(error, Context)

type Application struct {
	premiddleware []MiddlewareFunc
	middleware    []MiddlewareFunc
	Initialized   bool
	Debug         bool
	ErrorHandler  ErrorHandler
	Build         string
	router        *Router
}

var app *Application

func App() *Application {
	if app != nil {
		return app
	}
	app = &Application{
		ErrorHandler: app.DefaultErrorHandler,
		router:       NewRouter(app),
	}
	InitLogger()
	//app.Logger = log.New(ErrorLog, "", log.LstdFlags&log.Llongfile)

	return app
}

// Pre adds middleware to the chain which is run before router.
func (app *Application) Pre(middleware ...MiddlewareFunc) {
	app.premiddleware = append(app.premiddleware, middleware...)
}

// Use adds middleware to the chain which is run after router.
func (app *Application) Use(middleware ...MiddlewareFunc) {
	app.middleware = append(app.middleware, middleware...)
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (app *Application) Add(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.router.Add(path, func(c Context) error {
		h := applyMiddleware(handler, middleware...)
		return h(c)
	})
}

// Group creates a new router group with prefix and optional group-level middleware.
func (app *Application) Group(prefix string, m ...MiddlewareFunc) *Group {
	g := app.router.Group(prefix, m...)
	return g
}

func (app *Application) DefaultErrorHandler(err error, ctx Context) {
	he, ok := err.(*Error)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*Error); ok {
				he = herr
			}
		}
	} else {
		he = &Error{
			Code:    http.StatusInternalServerError,
			Errno:   http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	code := he.Code
	errno := he.Errno
	message := he.Message

	switch msg := he.Message.(type) {
	case string:
		if ctx.IsDebug() {
			message = server.Map{
				"code":  code,
				"errno": errno,
				"msg":   msg,
				"error": err.Error(),
				"rid":   ctx.RequestID(),
				"t":     time.Now().Unix(),
			}
		} else {
			message = server.Map{
				"code":  code,
				"errno": errno,
				"msg":   msg,
				"rid":   ctx.RequestID(),
				"t":     time.Now().Unix(),
			}
		}
	case json.Marshaler:
		// do nothing - this type knows how to format itself to JSON
	case error:
		message = server.Map{
			"code":  code,
			"errno": errno,
			"msg":   msg.Error(),
			"rid":   ctx.RequestID(),
			"t":     time.Now().Unix(),
		}
	}

	// Send response
	if ctx.ResponseWriter().Size == 0 {
		err = ctx.JSON(message)
		if err != nil {
			ctx.Logger().Error("%v", err)
		}
	} else {
		ctx.Logger().Error("%v", err)
	}
}

func (app *Application) Exec(r *Request, w *ResponseWriter) {
	c := &context{
		app:            app,
		responseWriter: w,
		request:        r,
	}
	if r == nil {
		c.Error(ErrBadRequest)
	}

	h := NotFoundHandler

	if app.premiddleware == nil {

		if handle, ok := app.router.Find(c.Path()); ok {
			h = applyMiddleware(handle, app.middleware...)
		}

	} else {
		h = func(c Context) error {
			handle, ok := app.router.Find(c.Path())
			if !ok {
				return ErrNotFound
			}
			handle = applyMiddleware(handle, app.middleware...)
			return handle(c)
		}
		h = applyMiddleware(h, app.premiddleware...)
	}

	if app.Initialized == false {
		c.Error(ErrServiceUnavailable)
		return
	}

	defer func() {
		rec := recover()
		if rec != nil {
			c.Logger().Warn("[msg: request fail] [err: %v]", rec)
			c.Error(fmt.Errorf("%v", rec))
		}
	}()

	// Execute chain
	if err := h(c); err != nil {
		c.Error(err)
	}
	return
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
