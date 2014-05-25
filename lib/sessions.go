package lib

import (
	"github.com/go4r/handy"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	DefaultSessionKeyPars              = [][]byte{[]byte(`2TZs4ESupxTybxm9JYXEeV6FuvI8YNKA`)}
	DefaultSessionStore sessions.Store = sessions.NewFilesystemStore("private/sessions", DefaultSessionKeyPars...)
	DefaultSessionName                 = "handy_session"

)

func init() {

	handy.Server.Context.SetProviderMap(handy.ContextProviderMap{

		"session-store" : func(c *handy.Context) func() interface{} {
			return func() interface{} {
				return DefaultSessionStore
			}
		},
		"session":func(c *handy.Context) func() interface{} {
			r := c.Get("request").(*http.Request)


			store := c.Get("session-store").(sessions.Store)
			name := c.Get("session-name").(string)

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
		"session-name":func(c *handy.Context) func() interface{} {
			return func() interface{} {
				return DefaultSessionName
			}
		},
	})
}


func Session(r interface{}) (*sessions.Session) {
	return handy.CContext(r).Get("session").(*sessions.Session)
}

func SessionSet(r, k, v interface{}) {
	Session(r).Values[k] = v
}


func SessionInt(r interface{}, k interface{}) int {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {

		if val, ok := val.(int); ok {
			return val
		}
	}
	return 0
}

func SessionString(r interface{}, k interface{}) string {
	sess := Session(r)
	if val, ok := sess.Values[k]; ok {
		if val, ok := val.(string); ok {
			return val
		}
	}
	return ""
}
