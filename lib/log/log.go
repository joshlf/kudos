package log

import (
	"fmt"
	"sync"

	"github.com/joshlf/kudos/lib/build"
)

type Level uint32

const (
	Debug Level = iota
	Verbose
	Info
	Warn
	Error
)

var names = []string{"Debug", "Verbose", "Info", "Warn", "Error"}

func (l Level) String() string {
	if l < Debug || l > Error {
		return fmt.Sprintf("UnknownLevel(%v)", uint32(l))
	}
	return names[l]
}

var defaultLevel = Info

type Printer interface {
	Print(a ...interface{})
	Printf(format string, a ...interface{})
	Println(a ...interface{})

	local()
}

type printer struct {
	logger *Logger
	level  Level
}

func (p printer) Print(a ...interface{})                 { p.logger.Print(p.level, a...) }
func (p printer) Printf(format string, a ...interface{}) { p.logger.Printf(p.level, format, a...) }
func (p printer) Println(a ...interface{})               { p.logger.Println(p.level, a...) }
func (p printer) local()                                 {}

// Logger is a leveled logger. It can be used either by
// calling one of the Print* methods and passing a level
// argument:
//  l.Println(Info, "Hello, World!")
// or by using the Printer fields:
//  l.Info.Println("Hello, World!")
//
// The zero value of a Logger is invalid - callers must
// use NewLogger.
type Logger struct {
	Debug, Verbose, Info, Warn, Error Printer
	level                             Level
	init                              bool
	m                                 sync.RWMutex
}

// NewLogger initializes and returns a new Logger
// whose logging level is set to the default (Debug
// if compiled in debug mode and Info otherwise).
func NewLogger() *Logger {
	l := &Logger{level: defaultLevel, init: true}
	l.Debug = printer{l, Debug}
	l.Verbose = printer{l, Verbose}
	l.Info = printer{l, Info}
	l.Warn = printer{l, Warn}
	l.Error = printer{l, Error}
	return l
}

// SetLevel sets the minimum level that will be printed.
// If a print method is called with a level less than
// the minimum (ie, a level m such that m < lvl), the
// method call is a no-op.
//
// The default is Info unless built in debug mode, in which
// case it is Debug.
func (l *Logger) SetLevel(lvl Level) {
	l.m.Lock()
	defer l.m.Unlock()
	l.checkInit()
	l.level = lvl
}

// GetLevel returns l's current logging level.
func (l *Logger) GetLevel() Level {
	l.m.Lock()
	defer l.m.Unlock()
	l.checkInit()
	return l.level
}

func (l *Logger) checkInit() {
	if !l.init {
		panic("lib/log: uninitialized Logger")
	}
}

func (l *Logger) Print(lvl Level, a ...interface{}) {
	l.m.RLock()
	defer l.m.RUnlock()
	l.checkInit()
	if lvl < l.level {
		return
	}
	fmt.Print(a...)
}

func (l *Logger) Printf(lvl Level, format string, a ...interface{}) {
	l.m.RLock()
	defer l.m.RUnlock()
	l.checkInit()
	if lvl < l.level {
		return
	}
	fmt.Printf(format, a...)
}

func (l *Logger) Println(lvl Level, a ...interface{}) {
	l.m.RLock()
	defer l.m.RUnlock()
	l.checkInit()
	if lvl < l.level {
		return
	}
	fmt.Println(a...)
}

var global *Logger

// SetLevel sets the level of the global logger to l.
func SetLevel(l Level) {
	global.SetLevel(l)
}

// Print calls Print on the global logger using the
// level l.
func (l Level) Print(a ...interface{}) {
	global.Print(l, a...)
}

// Printf calls Printf on the global logger using the
// level l.
func (l Level) Printf(format string, a ...interface{}) {
	global.Printf(l, format, a...)
}

// Println calls Println on the global logger using the
// level l.
func (l Level) Println(a ...interface{}) {
	global.Println(l, a...)
}

func init() {
	if build.DebugMode {
		defaultLevel = Debug
	}
	global = NewLogger()
}
