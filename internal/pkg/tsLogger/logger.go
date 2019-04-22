package tsLogger

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"fmt"
	"io"
	"log"
	"runtime/debug"
)

type logMessage struct {
	logger *log.Logger
	msg    string
}

// TODO: colorize prompt output
//const (
//	InfoColor    = "\033[1;34m%s\033[0m"
//	TraceColor   = "\033[1;36m%s\033[0m"
//	WarningColor = "\033[1;33m%s\033[0m"
//	ErrorColor   = "\033[1;31m%s\033[0m"
//	DebugColor   = "\033[0;36m%s\033[0m"
//)

var Logger = NewLogger()
var levels = config.Logger.Levels

type TSLogger struct {
	//traceChan   chan interface{}
	//infoChan    chan interface{}
	//warningChan chan interface{}
	//errorChan   chan interface{}
	//accChan     chan interface{}

	fatalChan chan interface{}

	logChan chan logMessage

	TraceLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	AccessLogger  *log.Logger

	FatalLogger *log.Logger
}

func NewLogger() *TSLogger {
	l := &TSLogger{
		//traceChan:   make(chan interface{}, 256),
		//infoChan:    make(chan interface{}, 256),
		//warningChan: make(chan interface{}, 256),
		//errorChan:   make(chan interface{}, 256),
		//accChan:     make(chan interface{}, 256),
		logChan:   make(chan logMessage, 256),
		fatalChan: make(chan interface{}, 256),
	}

	l.SetLoggers(
		levels["trace"],
		levels["info"],
		levels["warn"],
		levels["err"],
		levels["access"],
		levels["fatal"],
	)

	return l
}

func (l *TSLogger) LogTrace(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.logChan <- logMessage{l.TraceLogger, msg}
}

func (l *TSLogger) LogInfo(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.logChan <- logMessage{l.InfoLogger, msg}
}

func (l *TSLogger) LogWarn(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.logChan <- logMessage{l.WarningLogger, msg}
}

func (l *TSLogger) LogErr(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.logChan <- logMessage{l.ErrorLogger, msg}
}

func (l *TSLogger) LogAcc(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	l.logChan <- logMessage{l.AccessLogger, msg}
}

func (l *TSLogger) LogFatal(formatMessage string, values ...interface{}) {
	msg := fmt.Sprintf(formatMessage, values...)
	msg += "\n" + string(debug.Stack()) + "\n"
	l.fatalChan <- msg
}

func LogTrace(formatMessage string, values ...interface{}) {
	Logger.LogTrace(formatMessage, values...)
}

func LogInfo(formatMessage string, values ...interface{}) {
	Logger.LogInfo(formatMessage, values...)
}

func LogWarn(formatMessage string, values ...interface{}) {
	Logger.LogWarn(formatMessage, values...)
}

func LogErr(formatMessage string, values ...interface{}) {
	Logger.LogErr(formatMessage, values...)
}

func LogAcc(formatMessage string, values ...interface{}) {
	Logger.LogAcc(formatMessage, values...)
}

func LogFatal(formatMessage string, values ...interface{}) {
	Logger.LogFatal(formatMessage, values...)
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
			//case msg := <-l.traceChan:
			//	l.TraceLogger.Println(msg)
			//case msg := <-l.infoChan:
			//	l.InfoLogger.Println(msg)
			//case msg := <-l.warningChan:
			//	l.WarningLogger.Println(msg)
			//case msg := <-l.errorChan:
			//	l.ErrorLogger.Println(msg)
			//case msg := <-l.accChan:
			//	l.AccessLogger.Println(msg)
			case msg := <-l.logChan:
				msg.logger.Println(msg.msg)
			case msg := <-l.fatalChan:
				l.FatalLogger.Fatal(msg)
			}
		}
	}()
}

func (l *TSLogger) SetLoggers(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	accHandle io.Writer,
	fatalHandle io.Writer) {

	l.TraceLogger = log.New(traceHandle,
		"[TRACE]: ",
		log.Ldate|log.Ltime)

	l.InfoLogger = log.New(infoHandle,
		"[INFO]: ",
		log.Ldate|log.Ltime)

	l.WarningLogger = log.New(warningHandle,
		"[WARN]: ",
		log.Ldate|log.Ltime)

	l.ErrorLogger = log.New(errorHandle,
		"[ERROR]: ",
		log.Ldate|log.Ltime)

	l.AccessLogger = log.New(accHandle,
		"[ACC]: ",
		log.Ldate|log.Ltime)

	l.FatalLogger = log.New(fatalHandle,
		"[FATAL]: ",
		log.Ldate|log.Ltime)
}
