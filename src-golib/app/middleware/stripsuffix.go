package middleware

import (
	"strings"

	"github.com/lazygo/client/framework"
)

// StripUrlSuffix 去除url后缀
func StripUrlSuffix(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		path := ctx.Path()

		switch {
		case strings.HasSuffix(path, ".json"):
			path = strings.TrimSuffix(path, ".json")
		}

		ctx.Request().URL.Path = path

		return next(ctx)
	}
}
