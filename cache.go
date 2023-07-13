package router

import (
	"fmt"
	"net/http"
	"time"
)

func Cache(duration time.Duration) Middleware {
	return func(next HttpHandler) HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx *Context) {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, maxage=%v, s-maxage=%v, immutable", duration.Seconds(), duration.Seconds()))
			w.Header().Set("Expires", time.Now().Add(duration).Local().Format("Mon, 02 Jan 2006 15:04:05 MST"))
			next(w, r, ctx)
		}
	}
}
