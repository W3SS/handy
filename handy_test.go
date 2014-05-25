package handy_test

import (
	"github.com/go4r/handy"
	"net/http"
	"testing"
	"reflect"
)

type Get struct {
	Path string
}

type AnnotatedStruct struct {
	handy.Annotations
}

func (t *AnnotatedStruct) Annotates() {
	t.Annotation(
		Get{"/"},
		t.IndexAction,
	)
}

func (t *AnnotatedStruct) IndexAction(r *http.Request) {

}

func TestAnnotation(t *testing.T) {
	t.Log("Testing Annotations")
	var flag bool
	handy.GetAnnotations(&AnnotatedStruct{}).ProcessAnnotations(func(value interface{}, annotation interface{}) {
		switch annotation := annotation.(type) {
		case Get:
			if annotation.Path == "/" && value != nil {
				flag = true
			}
		}
	})

	if !flag {
		t.Fail()
	}
}

////////

var (
	Context = handy.NewContext()
)

type Session map[string]string

func TestContext(t *testing.T) {

	t.Log("Testing Context")

	Context.Set("session", func(c *handy.Context) func() interface{} {
		var cache = Session{}
		return func() interface{} {
			return cache
		}
	})

	session := Context.Map["session"]().(Session)
	session["userId"] = "1212"

	session2 := Context.Map["session"]().(Session)
	if session2["userId"] != "1212" {
		t.Fail()
	}
	t.Log("Testing Context 2")
	//Test 2
	var cache = Session{}
	Context.Set("session", func() interface{} {
		return cache
	})

	session = Context.Map["session"]().(Session)
	session["userId"] = "1212"

	session2 = Context.Map["session"]().(Session)
	if session2["userId"] != "1212" {
		t.Fail()
	}

}

//ParserParser Test

func TestParserParser(t *testing.T) {
	var parserParser handy.ParserParser
	parserParser.Build([]byte("/getusers/(:name)/(num:age)"))

	for _, v := range parserParser.Paths {
		switch v := v.(type){
		case []byte:
			t.Log(string(v))
		case *handy.ParameterParser:
			t.Logf("Name %s Parser %s\n", string(v.Name), string(v.Parser))
		}
	}

	var parameters = map[string]interface{}{}

	if parserParser.Test([]byte("/getusers/stylo/13"), nil, parameters) {
		if reflect.DeepEqual(parameters["name"], "stylo") == false {
			t.Errorf("%#v", parameters)
		}
		if reflect.DeepEqual(parameters["age"], "13") == false {
			t.Errorf("%#v", parameters)
		}
	}else {
		t.Fail()
	}
	if parserParser.Test([]byte("/getusersstylo/13"), nil, parameters) {
		t.Fail()
	}
}
