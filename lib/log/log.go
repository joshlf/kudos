package log

import (
	"fmt"
	"sync/atomic"
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

var level uint32

// SetLoggingLevel sets the minimum level
// that will be printed. For level m,
// m.Print* will execute only if m >= l;
// otherwise, it will be a no-op.
//
// The default is Info.
func SetLoggingLevel(l Level) {
	atomic.StoreUint32(&level, uint32(l))
}

func (l Level) Print(a ...interface{}) {
	if uint32(l) < atomic.LoadUint32(&level) {
		return
	}
	fmt.Print(a...)
}

func (l Level) Printf(format string, a ...interface{}) {
	if uint32(l) < atomic.LoadUint32(&level) {
		return
	}
	fmt.Printf(format, a...)
}

func (l Level) Println(a ...interface{}) {
	if uint32(l) < atomic.LoadUint32(&level) {
		return
	}
	fmt.Println(a...)
}

func init() {
	SetLoggingLevel(Info)
}
