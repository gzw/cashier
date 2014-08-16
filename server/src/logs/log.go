//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw

package logs

import (
	"fmt"
	"path"
	"runtime"
	"sync"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
)

type loggerType func() LoggerInterface

// LoggerInterface defines the behavior a log provider.
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name
// If Register is called twice with the same name or if driver is nil, it panics
func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider" + name)
	}
	adapters[name] = log
}

// DefaltLogger is default logger in cashier app
// it can contain several providers and log message into all providers
type DefaltLogger struct {
	lock                sync.Mutex
	level               int
	enableFunCallDepth  bool
	loggerFuncCallDepth int
	msg                 chan *logMsg
	outputs             map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

func NewLogger(channellen int64) *DefaltLogger {
	dl := new(DefaltLogger)
	dl.loggerFuncCallDepth = 2
	dl.msg = make(chan *logMsg, channellen)
	dl.outputs = make(map[string]LoggerInterface)
	go dl.startLogger()
	return dl
}

func (dl *DefaltLogger) SetLogger(adatername string, config string) error {
	dl.lock.Lock()
	defer dl.lock.Unlock()
	if log, ok := adapters[adatername]; ok {
		lg := log()
		err := lg.Init(config)
		dl.outputs[adatername] = lg
		if err != nil {
			fmt.Println("logs.DefaultLogger.SetLogger: " + err.Error())
			return err
		}
	} else {
		return fmt.Errorf("logs: unkown adaptname:%q(forgeten Register?)", adatername)
	}
	return nil
}

// remove a logger adapter in DefaltLogger
func (dl *DefaltLogger) DelLogger(adaptername string) error {
	dl.lock.Lock()
	defer dl.lock.Unlock()
	if lg, ok := dl.outputs[adaptername]; ok {
		lg.Destroy()
		delete(dl.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unkown adaptername %q(forgotten Register)", adaptername)
	}
}

func (dl *DefaltLogger) writeMsg(loglevel int, msg string) error {
	if dl.level > loglevel {
		return nil
	}
	lm := new(logMsg)
	lm.level = loglevel
	if dl.enableFunCallDepth {
		_, file, line, ok := runtime.Caller(dl.loggerFuncCallDepth)
		if ok {
			_, filename := path.Split(file)
			lm.msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	dl.msg <- lm
	return nil
}

func (dl *DefaltLogger) Setlevel(l int) {
	dl.level = l
}

func (dl *DefaltLogger) SetLogFuncCallDepth(d int) {
	dl.loggerFuncCallDepth = d
}

func (dl *DefaltLogger) EnableFuncCallDepth(b bool) {
	dl.enableFunCallDepth = b
}

func (dl *DefaltLogger) startLogger() {
	for {
		select {
		case dm := <-dl.msg:
			for _, l := range dl.outputs {
				l.WriteMsg(dm.msg, dm.level)
			}
		}
	}
}

func (dl *DefaltLogger) Close() {
	for {
		if len(dl.msg) > 0 {
			bm := <-dl.msg
			for _, l := range dl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		} else {
			break
		}
	}

	for _, l := range dl.outputs {
		l.Flush()
		l.Destroy()
	}
}

// log trace level message.
func (dl *DefaltLogger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	dl.writeMsg(LevelTrace, msg)
}

// log debug level message.
func (dl *DefaltLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	dl.writeMsg(LevelDebug, msg)
}

// log info level message.
func (dl *DefaltLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	dl.writeMsg(LevelInfo, msg)
}

// log warn level message.
func (dl *DefaltLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	dl.writeMsg(LevelWarn, msg)
}

// log error level message.
func (dl *DefaltLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	dl.writeMsg(LevelError, msg)
}

// log critical level message.
func (dl *DefaltLogger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	dl.writeMsg(LevelCritical, msg)
}
