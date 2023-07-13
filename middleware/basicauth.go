package middleware

import (
	"net/http"

	"github.com/gossie/router"
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

func BasicAuth(userChecker UserChecker) router.Middleware {
	return func(next router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, pathVariables map[string]string) {
			performBasicAuth(w, r, pathVariables, userChecker, next)
		}
	}
}

func performBasicAuth(w http.ResponseWriter, r *http.Request, pathVariables map[string]string, userChecker UserChecker, next router.HttpHandler) {
	if user, pass, ok := r.BasicAuth(); ok && userChecker(newUserData(user, pass)) {
		next(w, r, pathVariables)
		return
	}
	http.Error(w, "", http.StatusUnauthorized)
}
