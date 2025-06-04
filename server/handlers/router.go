package handlers

import (
	"habit-tracker/server/db"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type route struct {
	method  string
	pattern string
	handler HandlerFunc
}

type Router struct {
	routes []route
}

var Database db.Database

func CreateRouter() *Router {
	return &Router{
		routes: []route{},
	}
}

func (r *Router) Handle(method, pattern string, handler HandlerFunc) {
	r.routes = append(r.routes, route{method: method, pattern: pattern, handler: handler})
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Change to only allow requests from the frontend
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	addCORSHeaders(w)

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, route := range r.routes {
		if req.Method != route.method {
			continue
		}

		params, matched := Match(route.pattern, req.URL.Path)

		if !matched {
			continue
		}

		log.Println(route.method, req.URL.Path)

		route.handler(w, req, params)
		return
	}

	http.NotFound(w, req)
}

func Match(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i := range patternParts {
		if patternParts[i] == pathParts[i] {
			continue
		}

		if strings.HasPrefix(patternParts[i], ":") {
			params[strings.TrimPrefix(patternParts[i], ":")] = pathParts[i]
		} else {
			return nil, false
		}
	}

	return params, true
}
