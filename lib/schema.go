package lib

import (
	"net/http"
	"github.com/go4r/handy"
	"github.com/gorilla/schema"
)

var DefaultSchemaDecoder = schema.NewDecoder()

func init() {

	handy.Server.Context.SetProvider("gorilla.schema", func(c *handy.Context) func() interface{} {
		return func() interface{} {
			return DefaultSchemaDecoder
		}
	})

}


func SchemeDecoder(r interface{}) (*schema.Decoder) {
	return handy.CContext(r).Get("gorilla.schema").(*schema.Decoder)
}

func SchemeBindURL(r *http.Request, dest interface{}) error {
	if r.Form == nil {
		r.ParseForm()
	}
	return SchemeDecoder(r).Decode(dest, r.Form)
}

func SchemeBindBody(r *http.Request, dest interface{}) error {
	if r.PostForm == nil {
		r.ParseForm()
	}
	return SchemeDecoder(r).Decode(dest, r.PostForm)
}
