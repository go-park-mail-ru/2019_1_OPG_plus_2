package tsLogger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var Logger = NewLogger()

type TSLogger struct {
	traceChan   chan interface{}
	infoChan    chan interface{}
	warningChan chan interface{}
	errorChan   chan interface{}
	reqChan     chan interface{}

	TraceLogger        *log.Logger
	InfoLogger         *log.Logger
	WarningLogger      *log.Logger
	ErrorLogger        *log.Logger
	RequestBenchLogger *log.Logger
}

func NewLogger() *TSLogger {
	l := &TSLogger{
		traceChan:   make(chan interface{}, 256),
		infoChan:    make(chan interface{}, 256),
		warningChan: make(chan interface{}, 256),
		errorChan:   make(chan interface{}, 256),
		reqChan:     make(chan interface{}, 256),
	}
	l.InitLoggers(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stdout)
	return l
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func newStatusWriter(responseWriter http.ResponseWriter) *statusWriter {
	return &statusWriter{ResponseWriter: responseWriter}
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func (l *TSLogger) RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := newStatusWriter(w)
		next.ServeHTTP(sw, r)

		l.LogReq(fmt.Sprintf(
			"%v %s %s %s",
			sw.status,
			r.Method,
			r.RequestURI,
			time.Since(start),
		))
	})
}

func (l *TSLogger) LogTrace(msg interface{}) {
	l.traceChan <- msg
}

func (l *TSLogger) LogInfo(msg interface{}) {
	l.infoChan <- msg
}

func (l *TSLogger) LogWarn(msg interface{}) {
	l.warningChan <- msg
}

func (l *TSLogger) LogErr(msg interface{}) {
	l.errorChan <- msg
}

func (l *TSLogger) LogReq(msg interface{}) {
	l.reqChan <- msg
}

func (l *TSLogger) Run() {
	go func() {
		for {
			select {
			case msg := <-l.traceChan:
				l.TraceLogger.Println(msg)
			case msg := <-l.infoChan:
				l.InfoLogger.Println(msg)
			case msg := <-l.warningChan:
				l.WarningLogger.Println(msg)
			case msg := <-l.errorChan:
				l.ErrorLogger.Println(msg)
			case msg := <-l.reqChan:
				l.RequestBenchLogger.Println(msg)
			}

		}
	}()
}

func (l *TSLogger) InitLoggers(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	reqHandle io.Writer) {

	l.TraceLogger = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.InfoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.WarningLogger = log.New(warningHandle,
		"WARN: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.ErrorLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.RequestBenchLogger = log.New(reqHandle,
		"REQ: ",
		log.Ldate|log.Ltime)
}
