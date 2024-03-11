package router

import (
	"context"
	"net/http"
)

type usernamekey string

const UsernameKey = usernamekey("username")

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
		return func(w http.ResponseWriter, r *http.Request) {
			performBasicAuth(w, r, userChecker, next)
		}
	}
}

func performBasicAuth(w http.ResponseWriter, r *http.Request, userChecker UserChecker, next HttpHandler) {
	if user, pass, ok := r.BasicAuth(); ok && userChecker(newUserData(user, pass)) {

		next(w, r.WithContext(context.WithValue(r.Context(), UsernameKey, user)))
		return
	}
	http.Error(w, "", http.StatusUnauthorized)
}
