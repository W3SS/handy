package handy

import (
	"net/http"
	"bytes"
)

type responseWriter struct {
	http.ResponseWriter
	r *http.Request
	buffer *bytes.Buffer
	status int
}

func (w *responseWriter) Commit() {
	if w.buffer != nil {
		w.ResponseWriter.WriteHeader(w.status)

		w.ResponseWriter.Write(w.buffer.Bytes())
		w.buffer.Reset()
		w.buffer = nil
	}
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	//	w.ResponseWriter.WriteHeader(status)
	//	w.ResponseWriter.Write(w.buffer.Bytes())
	//	w.buffer.Reset()
	//	w.buffer = nil
}
func (w *responseWriter) Write(by []byte) (int, error) {
	if w.buffer == nil {
		return w.ResponseWriter.Write(by)
	}
	return w.buffer.Write(by)
}
