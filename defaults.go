package handy

import (
	"strconv"
	"errors"
	"net/http"
	"fmt"
)

var (
	ERROR_FORWARD_CONTROLLER_NOT_FOUND      = errors.New("Controller Not Found!")
	ERROR_FORWARD_WRONG_NUMBER_OF_ARGUMENTS = errors.New("Wrong Number Of Arguments")
	Server                                  = NewMux()
)

func ListenAndServer(addr string) error {
	return http.ListenAndServe(addr, Server)
}

func ListenAndServerTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(addr, cert, key, Server)
}


func Integer(url []byte) (interface{}, int, bool) {

	for k, v := range url {
		if v < '0' || v > '9' {
			if k == 0 {
				if v != '-' && v != '+' {
					return nil, 0, false
				}
			}else {
				value := string(url[0:k])
				i, _ := strconv.Atoi(value)
				return i, k, true
			}
		}
	}

	i, _ := strconv.Atoi(string(url))
	return i, len(url), true
}

func Number(url []byte) (interface{}, int, bool) {
	var p = -1
	for k, v := range url {
		if v < '0' || v > '9' {
			if v == '.' && p == -1 {
				p = k
			}else {

				if k == 0 {
					if v != '-' && v != '+' {
						return nil, 0, false
					}
				}else {
					value := string(url[0:k])
					i, _ := strconv.ParseFloat(value, 64)
					return i, k, true
				}
			}
		}
	}

	i, _ := strconv.ParseFloat(string(url), 64)
	return i, len(url), true
}

func WildCard(url []byte) (interface{}, int, bool) {
	return string(url), len(url), true
}

func NewMuxWithContext(c *Context) *Mux {
	return &Mux{
		context:c,
		subPatterns:map[string]func([]byte) (interface{}, int, bool){
			"*":WildCard,
			"int":Integer,
			"integer":Integer,
			"decimal":Integer,
			"number":Number,
			"float":Number,
		},
		wrappers:[]func(interface{}) func(mux *Mux, w http.ResponseWriter, r *http.Request, c *Context){
			defaultWrapper,
			jsonWrapper,
			textWrapper,
		},
	}
}

func NewMux() *Mux {
	return NewMuxWithContext(NewContext())
}

func Parameters(r interface{}) map[string]interface{} {
	c := CContext(r)
	if c == nil {
		return nil
	}

	v0, _ := c.Get("request.parameters").(map[string]interface{})

	return v0
}

func Parameter(r interface{}, key string) interface{} {
	return Parameters(r)[key]
}

func IntParameter(r interface{}, key string) int {
	return Parameter(r, key).(int)
}

func StringParameter(r interface{}, key string) string {
	return Parameter(r, key).(string)
}

func NumberParameter(r interface{}, key string) float64 {
	return Parameter(r, key).(float64)
}

func Forward(r interface{}, controller string, arguments ...interface{}) error {
	c := CContext(r)
	mux := c.Get("mux").(*Mux)
	writer := c.Get("response").(http.ResponseWriter)

	if parser, ok := mux.Names[controller]; ok {
		var (
			url string
			i      = 0
			length = len(arguments)
		)
		for _, v := range parser.Matcher.Parts {

			if i >= length {
				return ERROR_FORWARD_WRONG_NUMBER_OF_ARGUMENTS
			}

			switch v := v.(type){
			case []byte:
				url+=string(v)
			case *SubPatternMatcher:
				url += fmt.Sprintf("%v", arguments[i])
				i++
			}
		}

		if i < length-1 {
			return ERROR_FORWARD_WRONG_NUMBER_OF_ARGUMENTS
		}

		writer.Header().Set("Location", url)
		writer.WriteHeader(http.StatusFound)
		return nil
	}

	return ERROR_FORWARD_CONTROLLER_NOT_FOUND
}

func Map(controllers ...interface{}) *Mux {
	return Server.Map(controllers...)
}

func HandleAny(pattern string, handler interface{}) *MuxHandler {
	return Server.HandleAny(pattern, handler)
}

func HandleGet(pattern string, handler interface{}) *MuxHandler {
	return Server.HandleGet(pattern, handler)
}

func HandlePost(pattern string, handler interface{}) *MuxHandler {
	return Server.HandlePost(pattern, handler)
}

func HandlePut(pattern string, handler interface{}) *MuxHandler {
	return Server.HandlePut(pattern, handler)
}

func HandleDelete(pattern string, handler interface{}) *MuxHandler {
	return Server.HandleDelete(pattern, handler)
}

func HandleNotFound(pattern string, handlers interface{}) *MuxHandler {
	return Server.HandleNotFound(pattern, handlers)
}

func MiddlewareBefore(handlers ...interface{}) {
	Server.MiddlewareBefore(handlers...)
}

func MiddlewareAfter(handlers ...interface{}) {
	Server.MiddlewareAfter(handlers...)
}

func MiddlewareErrors(handlers ...interface{}) {
	Server.MiddlewareErrors(handlers...)
}
