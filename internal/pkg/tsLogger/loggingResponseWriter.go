package tsLogger

import "net/http"

type statusWriter struct {
	http.ResponseWriter
	Status        int
	ContentLength int
}

func NewStatusWriter(responseWriter http.ResponseWriter) *statusWriter {
	return &statusWriter{ResponseWriter: responseWriter}
}

func (w *statusWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.ContentLength += n
	return n, err
}
