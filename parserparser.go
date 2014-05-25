package handy

import (
	"net/http"
	"io"
	"errors"
)

const (
	INITIAL_STATE       = 0
	PARAM_PARSING_STATE = 1
	ESCAPING_STATE      = 1

)

type ParserParser struct {
	Paths  []interface{}
	Invoke func(mux *Mux, r * http.Request, context *Context)
}

func wrapper(handler interface{}) func(*Mux, * http.Request, *Context) {

	switch handler := handler.(type){

		//Handlers With Returns
	case func() string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler())
		}
	case func(*Context) string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler(context))
		}
	case func(*http.Request) string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler(r))
		}
	case func(http.ResponseWriter) string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler(w))
		}
	case func(http.ResponseWriter, *http.Request) string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler(w, r))
		}

	case func(*Context, *http.Request) string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler(context, r))
		}

	case func(*http.Request, *Context) string:
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			io.WriteString(w, handler(r, context))
		}

		//Handlers Without Returns
	case func(*Context):
		return func(mux *Mux, r *http.Request, context *Context) {
			handler(context)
		}
	case func(*Context, http.ResponseWriter):
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			handler(context, w)
		}
	case func(*Context, http.ResponseWriter, *http.Request):
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			handler(context, w, r)
		}
	case func(http.ResponseWriter, *http.Request):
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			handler(w, r)
		}
	case func(http.ResponseWriter):
		return func(mux *Mux, r *http.Request, context *Context) {
			w := context.Get("response").(http.ResponseWriter)
			handler(w)
		}
	case func(*http.Request):
		return func(mux *Mux, r *http.Request, context *Context) {
			handler(r)
		}
	case func(mux *Mux, r *http.Request, context *Context):
		return handler

	default:
		panic(errors.New("Invalid Handler"))
	}

	return nil
}


func MakeParserParser(pattern string, handler interface{}) *ParserParser {
	parser := &ParserParser{}
	parser.Build(([]byte)(pattern))
	parser.Invoke = wrapper(handler)
	return parser
}

func (parser *ParserParser) Try(url []byte, mux *Mux, r *http.Request, context * Context) bool {
	var parameters = map[string]interface{}{}

	if parser.Test(url, mux, parameters) {

		context.SetValue("parameters", parameters)

		parser.Invoke(mux, r, context)
		return true
	}

	return false
}

func (parser *ParserParser) Test(url []byte, mux *Mux, parameters map[string]interface{}) bool {
	var (
		i = 0
		l = len(url)
	)
	for _, v := range parser.Paths {
		switch v := v.(type){
		case []byte:
			//test here
		for k, v0 := range v {
			if i+k >= l {
				return false
			}
			if url[i+k] != v0 {
				return false
			}
		}
			i+=len(v)
		case *ParameterParser:
			var extractor func([]byte) (interface{}, int, bool)

			if mux != nil {
				if ex, ok := mux.ParametersParser[v.Parser]; ok {
					extractor = ex
				}
			}

			if extractor == nil {
				extractor = func(bytes []byte) (interface{}, int, bool) {
					for k, v0 := range bytes {
						if v0 == '/' {
							return string(bytes[0:k]), k, true
						}
					}
					return string(bytes), len(bytes), true
				}
			}

			value, numOfMatchedRunes, matchOk := extractor(url[i:])
			if matchOk == false {
				return false
			}

			if parameters != nil {
				parameters[v.Name] = value
			}

			i+=numOfMatchedRunes
		}
	}
	return len(url) == i
}

type ParameterParser struct {
	Name   string
	Parser string
}

func (parser *ParserParser) Build(s []byte) {
	var i int
	var l = len(s)
	var state = 0
	var parameterParser *ParameterParser

	for k := 0; k < l; k++ {
		var v = s[k]

		switch state {
		case PARAM_PARSING_STATE:

			if v == ')' {
				state = INITIAL_STATE
				parameterParser.Name = string(s[i:k])
				i = k+1
				parser.Paths = append(parser.Paths, parameterParser)
				continue
			}else if v == ':' {
				parameterParser.Parser = string(s[i:k])
				i = k+1
			}

			if k+1 >= l {
				panic("Invalid Pattern")
			}
		default:

			if v == '(' {
				state = PARAM_PARSING_STATE
				parser.Paths = append(parser.Paths, s[i:k])
				parameterParser = &ParameterParser{}
				i = k+1
			}
		}
	}

	if i < l {
		parser.Paths = append(parser.Paths, s[i:])
	}
}
