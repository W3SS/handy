package handy

import (
	"net/http"
	"bytes"
	"fmt"
)

type Any string
type Get string
type Post string
type Delete string
type Put string
type NotFound string

type Name string
type Before string
type After string


type Middleware func(r *http.Request) bool

type MuxControllerOnAttachMux interface {
	OnAttach(*Mux)
}

type MuxControllerOnAttachContext interface {
	OnAttach(*Context)
}

type MuxControllerOnRequest interface {
	OnRequest(r *http.Request) bool
}

type MuxControllerOnRequestEnd interface {
	OnRequestEnd(r *http.Request) bool
}

type MuxAnnotation interface {
	Attach(controller interface{}, action interface{}, mux *Mux)
}


type MuxController struct {
	Mux *Mux
	Value      interface{}
	before     []Middleware
	after      []Middleware
}

type MuxHandler struct {
	Controller *MuxController
	Matcher *PatternMatcher
	rawInvoker interface{}
	invoke func(*Mux, http.ResponseWriter, * http.Request, *Context)
	before     []Middleware
	after      []Middleware
}

func (handler *MuxHandler) Before(middlewares ...interface{}) *MuxHandler {
	for _, middleware := range middlewares {
		handler.before = append(handler.before, handler.Controller.Mux.GetMiddleware(middleware))
	}
	return handler
}

func (handler *MuxHandler) After(middlewares ...interface{}) *MuxHandler {
	for _, middleware := range middlewares {
		handler.after = append(handler.after, handler.Controller.Mux.GetMiddleware(middleware))
	}
	return handler
}

func (handler *MuxHandler) Invoke(mux *Mux, w http.ResponseWriter, r *http.Request, c *Context) {

	var isContinuable = true

	if handler.invoke == nil {
		handler.invoke = mux.Wrap(handler.rawInvoker)
	}



	if handler.Controller != nil {
		for _, middleware := range handler.Controller.before {
			if !middleware(r) {
				isContinuable = false
				break
			}
		}

		if isContinuable {
			if controller, ok := handler.Controller.Value.(MuxControllerOnRequest); ok {
				if !controller.OnRequest(r) {
					isContinuable = false
				}
			}
		}
	}

	if isContinuable {
		for _, middleware := range handler.before {
			if !middleware(r) {
				isContinuable = false
				break
			}
		}
	}
	if isContinuable {
		handler.invoke(mux, w, r, c)
	}

	for _, middleware := range handler.after {
		if !middleware(r) {
			return
		}
	}


	if handler.Controller != nil {

		if controller, ok := handler.Controller.Value.(MuxControllerOnRequestEnd); ok {
			if !controller.OnRequestEnd(r) {
				return
			}
		}

		for _, middleware := range handler.Controller.after {
			if !middleware(r) {
				return
			}
		}
	}
}

type Mux struct {
	Names    map[string]*MuxHandler

	any      []*MuxHandler
	get      []*MuxHandler
	post     []*MuxHandler
	put      []*MuxHandler
	delete   []*MuxHandler
	notFound []*MuxHandler

	preRequest      []Middleware
	postRequest     []Middleware
	errRequest      []Middleware


	wrappers    []func(interface{}) func(*Mux, http.ResponseWriter, * http.Request, *Context)

	subPatterns map[string]func([]byte) (interface{}, int, bool)

	context *Context
	defaultCtl MuxController
}

func (mux *Mux) SetDefaultCtl(ctl interface{}) {
	mux.defaultCtl.Value = ctl
}


func (mux *Mux) Wrap(handler interface{}) func(*Mux, http.ResponseWriter, * http.Request, *Context) {

	if handler, ok := handler.(func(*Mux, http.ResponseWriter, * http.Request, *Context)); ok {
		return handler
	}

	for _, wrapper := range mux.wrappers {
		handler := wrapper(handler)
		if handler != nil {
			return handler
		}
	}
	panic(fmt.Errorf("Invalid Handler %#v, Could not find a wrapper for the handler", handler))
}

func (mux *Mux) AddWrappers(wrapper func(interface{}) func(*Mux, http.ResponseWriter, * http.Request, *Context)) {
	mux.wrappers = append(mux.wrappers, wrapper)
}

func (mux *Mux) Context() *Context {
	return mux.context
}

func (mux *Mux) SetSubPattern(name string, matcher func([]byte) (interface{}, int, bool)) {
	if mux.subPatterns == nil {    mux.subPatterns = map[string]func([]byte) (interface{}, int, bool){}}
	mux.subPatterns[name] = matcher
}

func (mux *Mux) GetSubPattern(name string) func([]byte) (interface{}, int, bool) {

	if pattern, ok := mux.subPatterns[name]; ok {
		return pattern
	}

	if pattern, ok := mux.context.Get("mux.subPatterns." + name).(func([]byte) (interface{}, int, bool)); ok {
		mux.SetSubPattern(name, pattern)
		return pattern
	}

	return nil
}

func (mux *Mux) handlerLookup(url []byte, r *http.Request) (*MuxHandler, map[string]interface{}) {
	var muxHandlers []*MuxHandler

	switch r.Method{
	case "GET":
		muxHandlers = mux.get
	case "POST":
		muxHandlers = mux.post
	case "PUT":
		muxHandlers = mux.put
	case "DELETE":
		muxHandlers = mux.delete
	}

	for _, handler := range muxHandlers {
		var matched, parameters = handler.Matcher.Test(url, mux)
		if matched {
			return handler, parameters
		}
	}

	for _, handler := range mux.any {
		var matched, parameters = handler.Matcher.Test(url, mux)
		if matched {
			return handler, parameters
		}
	}

	for _, handler := range mux.notFound {
		var matched, parameters = handler.Matcher.Test(url, mux)
		if matched {
			return handler, parameters
		}
	}

	return nil, nil
}

func (mux *Mux) Dispatch(r_url string, w http.ResponseWriter, r *http.Request) {

	var url = []byte(r_url)

	context := beginRequestContext(r, mux.context)
	bufferedWriter := &responseWriter{w, r, bytes.NewBuffer(nil), http.StatusOK}


	context.SetValue("request", r)
	context.SetValue("response", bufferedWriter)
	context.SetValue("mux", mux)


	//Lookup for the matched handler
	var handler, parameters = mux.handlerLookup(url, r)
	context.SetValue("request.parameters", parameters)
	context.SetValue("request.handler", handler)

	defer func() {
		if err := recover(); err != nil {
			if mux.errRequest == nil {
				bufferedWriter.buffer.Reset()
				w.Write(([]byte)("Server Error"))
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}else {
				context.SetValue("panic.lasterror", err)
				for _, v := range mux.errRequest {
					if !v(r) {
						break
					}
				}
			}
			endRequestContext(r)
			bufferedWriter.Commit()
			return
		}
		endRequestContext(r)
		bufferedWriter.Commit()
	}()

	//Before
	for _, v := range mux.preRequest {
		if !v(r) {
			handler = nil //Disable the execution of the main handler
			break
		}
	}
	if handler != nil {
		handler.Invoke(mux, bufferedWriter, r, context)
	}
	//After
	for _, v := range mux.postRequest {
		if !v(r) {
			return
		}
	}
}


func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.Dispatch(r.URL.Path, w, r)
}

func (mux *Mux) Map(controllers ...interface{}) *Mux {

	for _, v := range controllers {

		var controller = &MuxController{Mux:mux, Value:v}

		switch v := v.(type){
		case MuxControllerOnAttachMux:
			v.OnAttach(mux)
		case MuxControllerOnAttachContext:
			v.OnAttach(mux.context)
		}

		for value, annotations := range GetAnnotations(v.(Annotated)) {

			if *value == v {
				for _, annotation := range annotations {
					switch annotation := annotation.(type){
					case Before:
						controller.before = append(controller.before, mux.GetMiddleware(string(annotation)))
					case After:
						controller.after = append(controller.after, mux.GetMiddleware(string(annotation)))
					}
				}
			}else {

				var name Name
				var parser *MuxHandler
				var before []Middleware
				var after []Middleware

				for _, annotation := range annotations {

					var slicePointer *[]*MuxHandler

					switch annotation := annotation.(type){
					case Name:
						name = annotation
					case Before:
						before = append(before, mux.GetMiddleware(string(annotation)))
					case After:
						after = append(after, mux.GetMiddleware(string(annotation)))
					case Get, Any, Post, Put, Delete:

						var pattern string
						switch annotation := annotation.(type){
						case Any:
							slicePointer = &mux.any
							pattern = string(annotation)
						case Get:
							slicePointer = &mux.get
							pattern = string(annotation)
						case Post:
							slicePointer = &mux.post
							pattern = string(annotation)
						case Put:
							slicePointer = &mux.put
							pattern = string(annotation)
						case Delete:
							slicePointer = &mux.delete
							pattern = string(annotation)
						case NotFound:
							slicePointer = &mux.notFound
							pattern = string(annotation)
						}

						parser = &MuxHandler{
							Matcher:MakePatternMatcher(pattern),
							before:before,
							after:after,
							rawInvoker:*value,
							Controller:controller,
						}
						if name != `` {
							if mux.Names == nil {
								mux.Names = map[string]*MuxHandler{}
							}
							mux.Names[(string)(name)] = parser
							name = ""
						}

						*slicePointer = append(*slicePointer, parser)

					case MuxAnnotation:
						annotation.Attach(v, *value, mux)
					}
				}
			}
		}
	}
	return mux
}

func (mux *Mux) MiddlewareBefore(arguments ...interface{}) {
	for _, v := range arguments {
		mux.preRequest = append(mux.preRequest, mux.GetMiddleware(v))
	}
}

func (mux *Mux) MiddlewareAfter(arguments ...interface{}) {
	for _, v := range arguments {
		mux.postRequest = append(mux.postRequest, mux.GetMiddleware(v))
	}
}

func (mux *Mux) MiddlewareErrors(arguments ...interface{}) {
	for _, v := range arguments {
		mux.errRequest = append(mux.errRequest, mux.GetMiddleware(v))
	}
}

func makeSimpleHandler(pattern string, handler interface{}, mux *Mux) *MuxHandler {
	if mux.defaultCtl.Mux == nil {
		mux.defaultCtl.Mux = mux
	}
	return &MuxHandler{Matcher:MakePatternMatcher(pattern), rawInvoker:handler, Controller:&mux.defaultCtl}
}

func (mux *Mux) HandleAny(pattern string, handler interface{}) *MuxHandler {
	value := makeSimpleHandler(pattern, handler, mux)
	mux.any = append(mux.any, value)
	return value
}

func (mux *Mux) HandleGet(pattern string, handler interface{}) *MuxHandler {
	value := makeSimpleHandler(pattern, handler, mux)
	mux.get = append(mux.get, value)
	return value
}

func (mux *Mux) HandlePost(pattern string, handler interface{}) *MuxHandler {
	value := makeSimpleHandler(pattern, handler, mux)
	mux.post = append(mux.post, value)
	return value
}

func (mux *Mux) HandlePut(pattern string, handler interface{}) *MuxHandler {
	value := makeSimpleHandler(pattern, handler, mux)
	mux.put = append(mux.put, value)
	return value
}

func (mux *Mux) HandleDelete(pattern string, handler interface{}) *MuxHandler {
	value := makeSimpleHandler(pattern, handler, mux)
	mux.delete = append(mux.delete, value)
	return value
}

func (mux *Mux) HandleNotFound(pattern string, handler interface{}) *MuxHandler {
	value := makeSimpleHandler(pattern, handler, mux)
	mux.notFound = append(mux.notFound, value)
	return value
}

