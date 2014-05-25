package handy

import (
	"net/http"
	"bytes"
)

type ResponseWriter struct {
	http.ResponseWriter
	r *http.Request
	buffer *bytes.Buffer
	status int
}

func (w *ResponseWriter) Commit() {
	if w.buffer != nil {
		w.ResponseWriter.WriteHeader(w.status)

		w.ResponseWriter.Write(w.buffer.Bytes())
		w.buffer.Reset()
		w.buffer = nil
	}
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
	//	w.ResponseWriter.WriteHeader(status)
	//	w.ResponseWriter.Write(w.buffer.Bytes())
	//	w.buffer.Reset()
	//	w.buffer = nil
}
func (w *ResponseWriter) Write(by []byte) (int, error) {
	if w.buffer == nil {
		return w.ResponseWriter.Write(by)
	}
	return w.buffer.Write(by)
}
