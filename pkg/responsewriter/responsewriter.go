package responsewriter

import (
	"net/http"
)

type StatusWriter struct {
	Writer http.ResponseWriter
	Status int
}

func (w *StatusWriter) Header() http.Header {
	return w.Writer.Header()
}

func (w *StatusWriter) WriteHeader(status int) {
	w.Status = status
	w.Writer.WriteHeader(status)
}

func (w *StatusWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = 200
	}
	n, err := w.Writer.Write(b)
	return n, err
}
