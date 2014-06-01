package handy

import (
	"net/http"
	"fmt"
)

func (mux *Mux) SetMiddleware(name string, middleware interface{}) {
	resolved := mux.GetMiddleware(middleware)
	mux.context.SetValue("middleware."+name, resolved)
}


func (mux *Mux) GetMiddleware(middleware interface{}) Middleware {

	switch middleware := middleware.(type){
	case Middleware:
		return middleware
	case string:
		return mux.GetMiddleware(mux.context.Get("middleware." + middleware))
	case func(*http.Request) bool:
		return middleware
	case func() bool:
		return func(*http.Request) bool {
			return middleware()
		}
	case func( *http.Request):
		return func(r*http.Request) bool {
			middleware(r)
			return true
		}
	case func():
		return func(*http.Request) bool {
			middleware()
			return true
		}
	case bool:
		return func(*http.Request) bool {
			return middleware
		}
	}

	panic(fmt.Errorf("Invalid Middleware %#v", middleware))

}
