package framework

import (
	stdContext "context"
	"io"
	"net/url"
)

type Request struct {
	ctx        stdContext.Context
	RequestURI string
	URL        *url.URL
	Data       io.Reader
}

func (r *Request) Context() stdContext.Context {
	return r.ctx
}

func NewRequestWithContext(ctx stdContext.Context, uri string, data io.Reader) (*Request, error) {
	u, err := url.Parse(uri)
	u.Path = cleanPath(u.Path)

	if err != nil {
		return nil, err
	}

	req := &Request{
		ctx:        ctx,
		RequestURI: uri,
		URL:        u,
		Data:       data,
	}
	return req, nil
}
