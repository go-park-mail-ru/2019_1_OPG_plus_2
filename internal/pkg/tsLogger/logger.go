package tsLogger

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"fmt"
	"io"
	"log"
)

var Logger = NewLogger()
var levels = config.Logger.Levels

type TSLogger struct {
	traceChan   chan interface{}
	infoChan    chan interface{}
	warningChan chan interface{}
	errorChan   chan interface{}
	accChan     chan interface{}

	TraceLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	AccessLogger  *log.Logger
}

func NewLogger() *TSLogger {
	l := &TSLogger{
		traceChan:   make(chan interface{}, 256),
		infoChan:    make(chan interface{}, 256),
		warningChan: make(chan interface{}, 256),
		errorChan:   make(chan interface{}, 256),
		accChan:     make(chan interface{}, 256),
	}
	l.SetLoggers(levels["trace"], levels["info"], levels["warn"], levels["err"], levels["access"])
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

func (l *TSLogger) LogAcc(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.accChan <- msg
}

func (l *TSLogger) Run() {
	go func() {
		defer func() {
			for _, f := range config.Logger.Files {
				err := f.Close()
				if err != nil {
					panic(err)
				}
			}
		}()

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
			case msg := <-l.accChan:
				l.AccessLogger.Println(msg)
			}
		}
	}()
}

func (l *TSLogger) SetLoggers(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	accHandle io.Writer) {

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

	l.AccessLogger = log.New(accHandle,
		"ACC: ",
		log.Ldate|log.Ltime)
}
