package framework

import (
	"io"
)

type ResponseWriter struct {
	Writer io.Writer
	Size   int64
}

// NewResponseWriter creates a new instance of Response.
func NewResponseWriter(w io.Writer) (r *ResponseWriter) {
	return &ResponseWriter{Writer: w}
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *ResponseWriter) Write(b []byte) (n int, err error) {
	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	return
}
