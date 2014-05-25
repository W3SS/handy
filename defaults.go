package handy

import (
	"strconv"
	"errors"
	"net/http"
	"fmt"
)

var (
	Server   = NewMux()
	dummyMap = map[string]interface{}{}
)

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
	myMux := &Mux{}

	myMux.ParametersParser = map[string]func([]byte) (interface{}, int, bool){
		"*":WildCard,
		"int":Integer,
		"integer":Integer,
		"decimal":Integer,
		"number":Number,
		"float":Number,
	}
	myMux.Names = map[string]*ParserParser{}
	myMux.Context = c

	return myMux
}

func NewMux() *Mux {
	return NewMuxWithContext(NewContext())
}

func Parameters(r interface{}) map[string]interface{} {
	c := CContext(r)
	if c == nil {
		return dummyMap
	}

	v0 := c.Get("parameters")
	if v0 == nil {
		return dummyMap
	}

	return v0.(map[string]interface{})
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

	if name, ok := mux.Names[controller]; ok {
		var url string
		var i = 0
		for _, v := range name.Paths {
			switch v := v.(type){
			case []byte:
				url+=string(v)
			case *ParameterParser:
				url += fmt.Sprintf("%v", arguments[i])
				i++
			}
		}

		writer.Header().Set("Location", url)
		writer.WriteHeader(http.StatusFound)
		//io.WriteString(writer, "Served")
		return nil
	}
	return errors.New("Controller Not Found!")
}


