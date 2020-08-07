package server

import (
	"context"
	"net/http"
	"regexp"
)

type ctxField struct{}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{
		method:  method,
		regex:   regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
	}
}

type router struct {
	routes []route
}

func newRouter(routes []route) *router {
	return &router{routes: routes}
}

func requestField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxField{}).([]string)
	return fields[index]
}

func (rt *router) handler(w http.ResponseWriter, r *http.Request) {
	for _, route := range rt.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			ctx := context.WithValue(r.Context(), ctxField{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}

	http.NotFound(w, r)
}
