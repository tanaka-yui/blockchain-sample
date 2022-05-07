package response

import (
	"fmt"
	"net/http"
)

type SendOptions struct {
	sMaxAge, maxAge int64
	OrgErr          error
}

type SendOption func(options *SendOptions)

func WithTTL(maxage, smaxage int64) SendOption {
	return func(o *SendOptions) {
		o.maxAge = maxage
		o.sMaxAge = smaxage
	}
}

func WithSMaxAge(v int64) SendOption {
	return func(o *SendOptions) {
		o.sMaxAge = v
	}
}

func WithError(err error) SendOption {
	return func(o *SendOptions) {
		o.OrgErr = err
	}
}

func SetCacheControl(h http.Header, opts *SendOptions) {
	const key = "Cache-Control"
	h.Set("Vary", "Accept,Accept-Encoding,Origin")
	if opts.maxAge > 0 {
		h.Add(key, fmt.Sprintf("max-age=%d", opts.maxAge))
	}
	if opts.sMaxAge > 0 {
		h.Add(key, fmt.Sprintf("s-maxage=%d", opts.sMaxAge))
	}
	if opts.sMaxAge == 0 && opts.maxAge == 0 {
		h.Set(key, "no-cache")
	}
}
