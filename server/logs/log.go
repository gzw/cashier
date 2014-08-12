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
	// log message levels
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
)

type loggerType func() LoggerInterface

// LoggerInterface defines the behavior of a log provider
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name
// If register is called twice with the same name or if driver is nil,
// it panics

func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: register provide is nil")
	}

	if _, dup := adapters[name]; dup {
		panic("logs: register called twice for provider" + name)
	}

	adapters[name] = log
}

type CashierLogger struct {
	lock                sync.Mutex
	level               int
	enableFuncCallDepth bool
	loggerFucnCallDepth int
	msg                 chan *logMsg
	outputs             map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

// NewLogger return a new CashierLogger.
// channellen means number of messages in chan
// if the buffering chan is full, logger adapters write to files or other way
func NewLogger(channellen int64) *CashierLogger {
	cl := new(CashierLogger)
	cl.loggerFucnCallDepth = 2
	cl.msg = make(chan *logMsg, channellen)
	cl.outputs = make(map[string]LoggerInterface)
	go cl.StartLogger()
	return cl
}

// SetLogger provides a given logger adapter into CashierLogger with config string
// config need to correct JSON as string : {"interval":360}
func (cl *CashierLogger) SetLogger(adaptername string, config string) error {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if log, ok := adapters[adaptername]; ok {
		lg := log()
		err := lg.Init(config)
		cl.outputs[adaptername] = lg
		if err != nil {
			fmt.Println("logs.CashierLogger.SetLogger:" + err.Error())
			return err
		}
	} else {
		return fmt.Errorf("logs: unknown adaptername %q(forgotten Register?)", adaptername)
	}
	return nil
}

// remove a logger adapter in CashierLogger.
func (cl *CashierLogger) DelLogger(adaptername string) error {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if lg, ok := cl.outputs[adaptername]; ok {
		lg.Destroy()
		delete(cl.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?", adaptername)
	}

}

func (cl *CashierLogger) WriteMsg(msg string, level int) error {
	if cl.level > level {
		return nil
	}
	lm := new(logMsg)
	lm.level = level
	if cl.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(cl.loggerFucnCallDepth)
		if ok {
			_, filename := path.Split(file)
			lm.msg = fmt.Sprintf("[%s:%d]%s", filename, line, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	cl.msg <- lm
	return nil
}

// set log message level.
// if message level (such as LevelTrace) is less than logger level (such as LevelWarn), ignore message
func (cl *CashierLogger) SetLevel(l int) {
	cl.level = l
}

// set log funcCallDepth
func (cl *CashierLogger) SetLogFuncCallDepth(d int) {
	cl.loggerFucnCallDepth = d
}

// enable log funcCallDepth
func (cl *CashierLogger) EnableFuncCallDepth(b bool) {
	cl.enableFuncCallDepth = b
}

func (cl *CashierLogger) StartLogger() {
	for {
		select {
		case cm := <-cl.msg:
			for _, l := range cl.outputs {
				l.WriteMsg(cm.msg, cm.level)
			}
		}
	}
}

// log trace level message
func (cl *CashierLogger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	cl.WriteMsg(msg, LevelTrace)
}

// log debug level message
func (cl *CashierLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	cl.WriteMsg(msg, LevelDebug)
}

// log info level message
func (cl *CashierLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	cl.WriteMsg(msg, LevelInfo)
}

// log error level message
func (cl *CashierLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	cl.WriteMsg(msg, LevelError)
}

// log critical level message
func (cl *CashierLogger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	cl.WriteMsg(msg, LevelCritical)
}

// flush all chan data
func (cl *CashierLogger) flush() {
	for _, l := range cl.outputs {
		l.Flush()
	}
}

// close logger, flush all chan data and destroy all adapters in cashierLogger
func (cl *CashierLogger) Close() {
	for {
		if len(cl.msg) > 0 {
			cm := <-cl.msg
			for _, l := range cl.outputs {
				l.WriteMsg(cm.msg, cm.level)
			}
		} else {
			break
		}
	}
	for _, l := range cl.outputs {
		l.Flush()
		l.Destroy()
	}
}
