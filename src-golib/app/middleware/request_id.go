package middleware

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/lazygo/client/framework"
)

// RequestID 添加request_id
func RequestID(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		rid := generator()
		ctx.WithValue(framework.HeaderXRequestID, rid)
		return next(ctx)
	}
}

func generator() uint64 {
	var x = strconv.Itoa(time.Now().Nanosecond() / 1000)
	res, errParseInt := strconv.ParseInt(x, 10, 64)
	if errParseInt != nil {
		return 0
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ((time.Now().Unix()*100000+res)&0xFFFFFFFF)*1000000000 + 100000000 + r.Int63n(899999999)
	return uint64(id)
}
