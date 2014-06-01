package handy

import (
	"net/http"
	"io"
	"encoding/json"
)

func textWrapper(handler interface{}) func(*Mux, http.ResponseWriter, * http.Request, *Context) {
	switch handler := handler.(type){
	case func() string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler())
		}
	case func(*Context) string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler(context))
		}
	case func(*http.Request) string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler(r))
		}
	case func(http.ResponseWriter) string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler(w))
		}
	case func(http.ResponseWriter, *http.Request) string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler(w, r))
		}
	case func(*Context, *http.Request) string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler(context, r))
		}
	case func(*http.Request, *Context) string:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			io.WriteString(w, handler(r, context))
		}
	}
	return nil
}


func jsonWrapper(handler interface{}) func(*Mux, http.ResponseWriter, * http.Request, *Context) {

	switch handler := handler.(type){

		//Handlers With String Returns

	case func() interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler())
		}
	case func(*Context) interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler(context))
		}
	case func(*http.Request) interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler(r))
		}
	case func(http.ResponseWriter) interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler(w))
		}
	case func(http.ResponseWriter, *http.Request) interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler(w, r))
		}
	case func(*Context, *http.Request) interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler(context, r))
		}
	case func(*http.Request, *Context) interface{}:
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(handler(r, context))
		}
	}
	return nil
}


func defaultWrapper(handler interface{}) func(*Mux, http.ResponseWriter, * http.Request, *Context) {

	switch handler := handler.(type){

	case func(*Context):
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			handler(context)
		}
	case func(*Context, http.ResponseWriter):
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			handler(context, w)
		}
	case func(*Context, http.ResponseWriter, *http.Request):
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			handler(context, w, r)
		}
	case func(http.ResponseWriter, *http.Request):
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			handler(w, r)
		}
	case func(http.ResponseWriter):
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			handler(w)
		}
	case func(*http.Request):
		return func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context) {
			handler(r)
		}
	case func(mux *Mux, w http.ResponseWriter, r *http.Request, context *Context):
		return handler
	}
	return nil
}
