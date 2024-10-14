package routing

import (
	"fmt"
	"net/http"
)

type HandlerFunc = http.HandlerFunc
type MiddlewareFunc = func(http.Handler) HandlerFunc

type RouterHandler struct {
	handler http.Handler
}

func New() *RouterHandler {
	return &RouterHandler{
		handler: http.NewServeMux(),
	}
}

func (r *RouterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r *RouterHandler) UseMiddleware(f MiddlewareFunc) *RouterHandler {
	r.handler = f(r.handler)
	return r
}

func (r *RouterHandler) Handle(pattern string, handler HandlerFunc) *RouterHandler {
	r.handler.(*http.ServeMux).Handle(pattern, handler)
	return r
}

func (r *RouterHandler) Get(pattern string, handler HandlerFunc) *RouterHandler {
	r.Handle(fmt.Sprintf("GET %s", pattern), handler)
	return r
}

func (r *RouterHandler) Post(pattern string, handler HandlerFunc) *RouterHandler {
	r.Handle(fmt.Sprintf("POST %s", pattern), handler)
	return r
}
func (r *RouterHandler) Delete(pattern string, handler HandlerFunc) *RouterHandler {
	r.Handle(fmt.Sprintf("DELETE %s", pattern), handler)
	return r
}
