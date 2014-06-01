package lib

import (
	"github.com/go4r/handy"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
)

var (
	DefaultSessionKeyPars              = [][]byte{[]byte(`2TZs4ESupxTybxm9JYXEeV6FuvI8YNKA`)}
	DefaultSessionStore sessions.Store = sessions.NewFilesystemStore("private/sessions", DefaultSessionKeyPars...)
	DefaultSessionName                 = "handy_session"
)

func init() {
	handy.Server.Context().MapProviders(handy.ProvidersMap{
		"cookies":func(c *handy.Context) func() interface{} {
			var cookies = &Cookies{map[string]*http.Cookie{}, c.Get("request").(*http.Request)}
			c.CleanupFunc(func() {
				w := c.Get("response").(http.ResponseWriter)
				for _, v := range cookies.Cookies {
					v.Value = url.QueryEscape(v.Value)
					http.SetCookie(w, v)
				}
			})
			return func() interface{} {
				return cookies
			}
		},
		"session.store" : func(c *handy.Context) func() interface{} {
			return func() interface{} {
				return DefaultSessionStore
			}
		},
		"session":func(c *handy.Context) func() interface{} {
			r := c.Get("request").(*http.Request)


			store := c.Get("session.store").(sessions.Store)
			name := c.Get("session.name").(string)

			value, err := store.Get(r, name)

			c.CleanupFunc(func() {
				w := c.Get("response").(http.ResponseWriter)
				value.Save(r, w)
			})

			return func() interface{} {
				if err != nil {
					return value
				}
				return value
			}
		},
		"session.name":func(c *handy.Context) func() interface{} {
			return func() interface{} {
				return DefaultSessionName
			}
		},
	})
}

type Cookies struct {
	Cookies map[string]*http.Cookie
	r *http.Request
}

func (cookies *Cookies) Set(k string, v string) *http.Cookie {
	if cookies.Cookies == nil {
		cookies.Cookies = map[string]*http.Cookie{}
	}
	cookie := &http.Cookie{Name:k, Value:v}
	cookies.Cookies[k] = cookie
	return cookie
}

func (cookies *Cookies) Get(k string) string {
	if v, ok := cookies.Cookies[k]; ok {
		return v.Value
	}
	value, _ := cookies.r.Cookie(k)
	return value.Value
}

func (cookies *Cookies) Has(k string) bool {
	if _, ok := cookies.Cookies[k]; ok {
		return ok
	}
	_, ok := cookies.r.Cookie(k)
	return ok == nil
}

func GetCookies(r interface{}) (*Cookies) {
	return handy.CContext(r).Get("cookies").(*Cookies)
}

func SetCookie(r interface{}, k, v string) *http.Cookie {
	return handy.CContext(r).Get("cookies").(*Cookies).Set(k, v)
}

func GetCookie(r interface{}, k string) string {
	return handy.CContext(r).Get("cookies").(*Cookies).Get(k)
}

func HasCookie(r interface{}, k string) bool {
	return handy.CContext(r).Get("cookies").(*Cookies).Has(k)
}

func Session(r interface{}) (*sessions.Session) {
	return handy.CContext(r).Get("session").(*sessions.Session)
}

func SetSession(r, k, v interface{}) {
	Session(r).Values[k] = v
}

func GetSession(r, k interface{}, v ...interface{}) interface{} {
	sess := Session(r).Values

	if value, ok := sess[k]; ok {
		return value
	}

	if len(v) != 0 {
		return v[0]
	}

	return nil
}


func HasSession(r interface{}, k interface{}) bool {
	sess := Session(r)
	_ , ok := sess.Values[k]
	return ok
}


func IntSession(r interface{}, k interface{}) int {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {

		if val, ok := val.(int); ok {
			return val
		}
	}
	return 0
}

func StringSession(r interface{}, k interface{}) string {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {
		if val, ok := val.(string); ok {
			return val
		}
	}
	return ""
}

func Float64Session(r interface{}, k interface{}) float64 {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {
		if val, ok := val.(float64); ok {
			return val
		}
	}
	return 0
}

func FloatSession(r interface{}, k interface{}) float32 {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {
		if val, ok := val.(float32); ok {
			return val
		}
	}
	return 0
}

func BoolSession(r interface{}, k interface{}) bool {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {
		if val, ok := val.(bool); ok {
			return val
		}
	}
	return false
}
