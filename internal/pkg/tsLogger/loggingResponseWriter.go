package tsLogger

import "net/http"

type StatusWriter struct {
	http.ResponseWriter
	Status        int
	ContentLength int
}

func NewStatusWriter(responseWriter http.ResponseWriter) *StatusWriter {
	return &StatusWriter{ResponseWriter: responseWriter}
}

func (w *StatusWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *StatusWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = 200
		//w.WriteHeader(200)
	}
	n, err := w.ResponseWriter.Write(b)
	w.ContentLength += n
	return n, err
}

func (w *StatusWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
