package handy

import (
	"io"
	"net/http"
	"sync"
	"errors"
)

var (
	globalRequestMapMutex sync.RWMutex
	globalRequestMap = map[*http.Request]*Context{}
)

func beginRequestContext(r *http.Request, c *Context) *Context {
	globalRequestMapMutex.Lock()

	if value, ok := globalRequestMap[r]; ok {
		globalRequestMapMutex.Unlock()
		return value
	}

	value := c.NewContext()
	globalRequestMap[r] = value
	globalRequestMapMutex.Unlock()
	return value
}

func endRequestContext(r *http.Request) {
	globalRequestMapMutex.Lock()
	globalRequestMap[r].GC()
	delete(globalRequestMap, r)
	globalRequestMapMutex.Unlock()
}

func Handler(invoker interface{}, c *Context) func(http.ResponseWriter, *http.Request) {
	switch invoker := invoker.(type) {
	case func(http.ResponseWriter):
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				endRequestContext(r)
			}()
			beginRequestContext(r, c)
			invoker(w)
		}
	case func(*http.Request) string:
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				endRequestContext(r)
			}()
			beginRequestContext(r, c)
			io.WriteString(w, invoker(r))
		}
	case func(http.ResponseWriter, *http.Request):
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				endRequestContext(r)
			}()
			beginRequestContext(r, c)
			invoker(w, r)
		}
	case func() string:
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				endRequestContext(r)
			}()
			beginRequestContext(r, c)
			io.WriteString(w, invoker())
		}
	}
	panic(errors.New("Invalid Handler"))
	return nil
}

func CContext(r interface{}) *Context {
	switch r := r.(type){
	case *http.Request:
		globalRequestMapMutex.RLock()
		if value, ok := globalRequestMap[r]; ok {
			globalRequestMapMutex.RUnlock()
			return value
		}
		globalRequestMapMutex.RUnlock()
		return nil
	case *Mux:
		return r.Context
	case *Context:
		return r
	}
	return nil
}
