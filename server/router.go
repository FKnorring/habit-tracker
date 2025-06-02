package main

import (
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

func CreateRouter() *Router {
	return &Router{
		routes: []route{},
	}
}

func (r *Router) Handle(method, pattern string, handler HandlerFunc) {
	r.routes = append(r.routes, route{method: method, pattern: pattern, handler: handler})
}

// addCORSHeaders adds the necessary CORS headers to the response
func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Add CORS headers to all responses
	addCORSHeaders(w)

	// Handle preflight OPTIONS requests
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, route := range r.routes {

		if req.Method != route.method {
			continue
		}

		params, matched := match(route.pattern, req.URL.Path)

		if !matched {
			continue
		}

		log.Println(route.method, req.URL.Path)

		route.handler(w, req, params)
		return
	}

	http.NotFound(w, req)
}

func match(pattern, path string) (map[string]string, bool) {
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
