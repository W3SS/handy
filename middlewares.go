package handy

import (
	"net/http"
	"fmt"
)

func (mux *Mux) SetMiddleware(name string, middleware interface{}) *Mux {
	resolved := mux.WrapMiddleware(middleware)
	mux.context.SetValue("middleware."+name, resolved)
	return mux
}


func (mux *Mux) WrapMiddleware(middleware interface{}) Middleware {

	switch middleware := middleware.(type){

	case string:
		return mux.WrapMiddleware(mux.context.Get("middleware." + middleware))

	case Middleware:
		return middleware
	case func(http.ResponseWriter, *http.Request, *Context) bool:
		return middleware

		//Just With Request
	case func(*http.Request) bool:
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			return middleware(r)
		}

	case func( *http.Request):
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			middleware(r)
			return true
		}

		//Just With Context
	case func(*Context) bool:
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			return middleware(c)
		}

	case func( *Context):
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			middleware(c)
			return true
		}

		//Just With ResponseWriter
	case func(http.ResponseWriter) bool:
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			return middleware(w)
		}

	case func(http.ResponseWriter):
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			middleware(w)
			return true
		}


	case func() bool:
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			return middleware()
		}
	case func():
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			middleware()
			return true
		}
	case bool:
		return func(w http.ResponseWriter, r *http.Request, c *Context) bool {
			return middleware
		}
	}
	panic(fmt.Errorf("Invalid Middleware %#v", middleware))
}
