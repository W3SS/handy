package handy

import (
	"net/http"
	"io"
	"encoding/json"
)

func CResponseWriter(r interface{}) http.ResponseWriter {
	if w, is := r.(http.ResponseWriter); is {
		return w
	}
	return CContext(r).Get("response").(http.ResponseWriter)
}

func SetContentType(r interface{}, contentType string) http.ResponseWriter {
	w := CResponseWriter(r)
	w.Header().Set("Content-Type", contentType)
	return w
}


func Respond(r interface{}, out interface{}) error {
	w := CResponseWriter(r)
	switch out := out.(type){
	case string:
		_, err := io.WriteString(w, out)
		return err
	default:
		return json.NewEncoder(w).Encode(r)
	}
}

func RespondWithStatus(r interface{}, out interface{}, status int) error {
	w := CResponseWriter(r)
	w.WriteHeader(status)
	return Respond(w, out)
}
