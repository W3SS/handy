package handy

import (
	"net/http"
	"sync"
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
	case Contextualized:
		return r.Context()
	case *Context:
		return r
	}
	return nil
}
