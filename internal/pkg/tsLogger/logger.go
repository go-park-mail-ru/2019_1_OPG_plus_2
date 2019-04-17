package tsLogger

import (
	"fmt"
	"io"
	"log"
	"os"
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
	l.SetLoggers(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stdout)
	return l
}

func (l *TSLogger) LogTrace(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.traceChan <- msg
}

func (l *TSLogger) LogInfo(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.infoChan <- msg
}

func (l *TSLogger) LogWarn(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.warningChan <- msg
}

func (l *TSLogger) LogErr(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.errorChan <- msg
}

func (l *TSLogger) LogReq(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
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

func (l *TSLogger) SetLoggers(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	reqHandle io.Writer) {

	l.TraceLogger = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime)

	l.InfoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	l.WarningLogger = log.New(warningHandle,
		"WARN: ",
		log.Ldate|log.Ltime)

	l.ErrorLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)

	l.RequestBenchLogger = log.New(reqHandle,
		"REQ: ",
		log.Ldate|log.Ltime)
}
