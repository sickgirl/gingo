package http

import (
	"net/http"

	"github.com/douglarek/zerodown"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HandleFunc func(r *Request) *Response
type HandleExtendFunc func(r *Request) (*Request, error)

type Router struct {
	*mux.Router
	corsOpts []handlers.CORSOption
}

type Route struct {
	*mux.Route
}

type OriginValidator func(string) bool

func NewRouter() *Router {
	return &Router{
		mux.NewRouter(),
		make([]handlers.CORSOption, 0),
	}
}

func NewSubRouter(path string) *Router {
	router := mux.NewRouter().PathPrefix(path).Subrouter()
	return &Router{
		router,
		make([]handlers.CORSOption, 0),
	}
}

func (r *Router) RouteHandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return &Route{r.HandleFunc(path, f)}
}

func (r *Router) AllowedMethods(methods []string) {
	r.corsOpts = append(r.corsOpts, handlers.AllowedMethods(methods))
}

func (r *Router) AllowedHeaders(headers []string) {
	r.corsOpts = append(r.corsOpts, handlers.AllowedHeaders(headers))
}

func (r *Router) AllowedOrigins(origins []string) {
	r.corsOpts = append(r.corsOpts, handlers.AllowedOrigins(origins))
}

func (r *Router) AllowedOriginValidator(fn OriginValidator) {
	r.corsOpts = append(r.corsOpts, handlers.AllowedOriginValidator(handlers.OriginValidator(fn)))
}

func (r *Router) ExposedHeaders(headers []string) {
	r.corsOpts = append(r.corsOpts, handlers.ExposedHeaders(headers))
}

func (r *Router) MaxAge(age int) {
	r.corsOpts = append(r.corsOpts, handlers.MaxAge(age))
}

func (r *Router) IgnoreOptions() {
	r.corsOpts = append(r.corsOpts, handlers.IgnoreOptions())
}

func (r *Router) AllowCredentials() {
	r.corsOpts = append(r.corsOpts, handlers.AllowCredentials())
}

// ListenAndServe This function blocks
func (r *Router) ListenAndServe(addr string) error {
	s := http.NewServeMux()
	s.Handle("/", handlers.HTTPMethodOverrideHandler(r))
	return zerodown.ListenAndServe(addr, s)
}

// ListenAndServeCORS This function blocks
func (r *Router) ListenAndServeCORS(addr string) error {
	s := http.NewServeMux()
	s.Handle("/", handlers.HTTPMethodOverrideHandler(r))
	return zerodown.ListenAndServe(addr, handlers.CORS(r.corsOpts...)(r))
}
