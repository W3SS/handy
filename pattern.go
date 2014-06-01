package handy

import "fmt"

const (
	INITIAL_STATE       = 0
	PARAM_PARSING_STATE = 1
)

type PatternMatcher struct {
	Parts  []interface{}
}

func MakePatternMatcher(pattern string) *PatternMatcher {
	parser := &PatternMatcher{}
	parser.Build(([]byte)(pattern))
	return parser
}

func (parser *PatternMatcher) Gen(arguments ...interface{}) string {
	var url string
	var i = 0
	for _, v := range parser.Parts {
		switch v := v.(type){
		case []byte:
			url+=string(v)
		case *SubPatternMatcher:
			url += fmt.Sprintf("%v", arguments[i])
			i++
		}
	}
	return url
}


func (parser *PatternMatcher) Test(url []byte, mux *Mux) (bool, map[string]interface{}) {
	var (
		i          = 0
		l          = len(url)
		parameters = map[string]interface{}{}
	)
	for _, v := range parser.Parts {
		switch v := v.(type){
		case []byte:
			//test here
		for k, v0 := range v {
			if i+k >= l {
				return false, nil
			}
			if url[i+k] != v0 {
				return false, nil
			}
		}
			i+=len(v)
		case *SubPatternMatcher:
			var subPatternMatcher func([]byte) (interface{}, int, bool)

			if mux != nil {
				subPatternMatcher = mux.GetSubPattern(v.Parser)
			}

			if subPatternMatcher == nil {
				subPatternMatcher = func(bytes []byte) (interface{}, int, bool) {
					for k, v0 := range bytes {
						if v0 == '/' {
							return string(bytes[0:k]), k, true
						}
					}
					return string(bytes), len(bytes), true
				}
			}

			value, numOfMatchedRunes, matchOk := subPatternMatcher(url[i:])
			if matchOk == false {
				return false, nil
			}

			if parameters != nil {
				parameters[v.Name] = value
			}

			i+=numOfMatchedRunes
		}
	}
	return len(url) == i, parameters
}

type SubPatternMatcher struct {
	Name   string
	Parser string
}

func (parser *PatternMatcher) Build(s []byte) {
	var i int
	var l = len(s)
	var state = 0
	var parameterParser *SubPatternMatcher

	for k := 0; k < l; k++ {
		var v = s[k]

		switch state {
		case PARAM_PARSING_STATE:

			if v == ')' {
				state = INITIAL_STATE
				parameterParser.Name = string(s[i:k])
				i = k+1
				parser.Parts = append(parser.Parts, parameterParser)
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
				parser.Parts = append(parser.Parts, s[i:k])
				parameterParser = &SubPatternMatcher{}
				i = k+1
			}
		}
	}

	if i < l {
		parser.Parts = append(parser.Parts, s[i:])
	}
}
