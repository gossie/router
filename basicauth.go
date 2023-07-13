package router

import (
	"net/http"
)

type UserData struct {
	username, password string
}

func newUserData(username, password string) *UserData {
	return &UserData{username, password}
}

func (ud *UserData) Username() string {
	return ud.username
}

func (ud *UserData) Password() string {
	return ud.password
}

type UserChecker = func(*UserData) bool

func BasicAuth(userChecker UserChecker) Middleware {
	return func(next HttpHandler) HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx *Context) {
			performBasicAuth(w, r, ctx, userChecker, next)
		}
	}
}

func performBasicAuth(w http.ResponseWriter, r *http.Request, ctx *Context, userChecker UserChecker, next HttpHandler) {
	if user, pass, ok := r.BasicAuth(); ok && userChecker(newUserData(user, pass)) {
		ctx.username = user
		next(w, r, ctx)
		return
	}
	http.Error(w, "", http.StatusUnauthorized)
}
