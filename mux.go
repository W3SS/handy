package handy

import (
	"net/http"
	"bytes"
)

type Any string
type Get string
type Post string
type Delete string
type Put string
type Name string

type Mux struct {
	Names    map[string]*ParserParser

	ANY    []*ParserParser
	GET    []*ParserParser
	POST   []*ParserParser
	PUT    []*ParserParser
	DELETE []*ParserParser

	preRequest      []func(*Mux, *http.Request, *Context)
	postRequest     []func(*Mux, *http.Request, *Context)
	notFoundRequest []func(*Mux, *http.Request, *Context)
	errRequest      []func(*Context, interface{})

	ParametersParser map[string]func([]byte) (interface{}, int, bool)

	Context *Context
}

func (mux *Mux) Dispatch(r_url string, w http.ResponseWriter, r *http.Request) {
	var url = []byte(r_url)
	context := beginRequestContext(r, mux.Context)

	context.SetValue("request", r)

	bufferedWriter := &ResponseWriter{w, r, bytes.NewBuffer(nil), http.StatusOK}

	context.SetValue("response", bufferedWriter)
	context.SetValue("mux", mux)


	defer func() {

		if err := recover(); err != nil {
			endRequestContext(r)
			bufferedWriter.Commit()

			if mux.errRequest == nil {
				//w.WriteHeader(http.StatusInternalServerError)
				//w.Write(([]byte)("Unexpected Error"))
				panic(err)
			}else {
				for _, v := range mux.errRequest {
					v(context, err)
				}
			}

			return
		}

		for _, v := range mux.postRequest {
			v(mux, r, context)
		}

		endRequestContext(r)
		bufferedWriter.Commit()
	}()

	for _, v := range mux.preRequest {
		v(mux, r, context)
	}
	//

	switch r.Method{
	case "GET":
	for _, v := range mux.GET {
		if v.Try(url, mux, r, context) {
			return
		}
	}

	case "POST":
	for _, v := range mux.POST {
		if v.Try(url, mux, r, context) {
			return
		}
	}

	case "PUT":
	for _, v := range mux.PUT {
		if v.Try(url, mux, r, context) {
			return
		}
	}

	case "DELETE":
	for _, v := range mux.DELETE {
		if v.Try(url, mux, r, context) {
			return
		}
	}
	}

	for _, v := range mux.ANY {
		if v.Try(url, mux, r, context) {
			return
		}
	}

	for _, v := range mux.notFoundRequest {
		v(mux, r, context)
	}
}


func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.Dispatch(r.URL.Path, w, r)
}

func (mux *Mux) Map(controllers ...interface{}) *Mux {
	for _, v := range controllers {

		var name Name

		GetAnnotations(v.(Annotated)).ProcessAnnotations(func(value interface{}, annotation interface{}) {
			var slicePointer *[]*ParserParser
			var parser *ParserParser
			switch annotation := annotation.(type){
			case Name:
				name = annotation

			case Get, Any, Post, Put, Delete:

				switch annotation := annotation.(type){
				case Any:
					parser = MakeParserParser((string)(annotation), value)
					slicePointer = &mux.ANY

				case Get:
					parser = MakeParserParser((string)(annotation), value)
					slicePointer = &mux.GET

				case Post:
					parser = MakeParserParser((string)(annotation), value)
					slicePointer = &mux.POST

				case Put:
					parser = MakeParserParser((string)(annotation), value)
					slicePointer = &mux.PUT

				case Delete:
					parser = MakeParserParser((string)(annotation), value)
					slicePointer = &mux.DELETE
				}

				if name != `` {
					mux.Names[(string)(name)] = parser
				}

				*slicePointer = append(*slicePointer, parser)
			}
		})

	}
	return mux
}

func (mux *Mux) Before(arguments ...interface{}) {
	for _, v := range arguments {
		mux.preRequest = append(mux.preRequest, wrapper(v))
	}
}

func (mux *Mux) After(arguments ...interface{}) {
	for _, v := range arguments {
		mux.postRequest = append(mux.postRequest, wrapper(v))
	}
}

func (mux *Mux) NotFound(arguments ...interface{}) {
	for _, v := range arguments {
		mux.notFoundRequest = append(mux.notFoundRequest, wrapper(v))
	}
}

func (mux *Mux) Errors(arguments ...func(*Context, interface{})) {
	mux.errRequest = append(mux.errRequest, arguments...)
}


func (mux *Mux) Any(pattern string, handler interface{}) {
	mux.ANY = append(mux.ANY, MakeParserParser(pattern, handler))
}

func (mux *Mux) Get(pattern string, handler interface{}) {
	mux.GET = append(mux.GET, MakeParserParser(pattern, handler))
}

func (mux *Mux) Post(pattern string, handler interface{}) {
	mux.POST = append(mux.POST, MakeParserParser(pattern, handler))
}

func (mux *Mux) Put(pattern string, handler interface{}) {
	mux.PUT = append(mux.PUT, MakeParserParser(pattern, handler))
}

func (mux *Mux) Delete(pattern string, handler interface{}) {
	mux.DELETE = append(mux.DELETE, MakeParserParser(pattern, handler))
}
