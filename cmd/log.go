package main

import "fmt"

type Level int

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
		return fmt.Sprintf("UnknownLevel(%v)", int(l))
	}
	return names[l]
}

var level Level

func SetLoggingLevel(l Level) {
	level = l
}

func (l Level) Print(a ...interface{}) {
	if l < level {
		return
	}
	fmt.Print(a...)
}

func (l Level) Printf(format string, a ...interface{}) {
	if l < level {
		return
	}
	fmt.Printf(format, a...)
}

func (l Level) Println(a ...interface{}) {
	if l < level {
		return
	}
	fmt.Println(a...)
}

func (l Level) prefix() string {
	if l == Debug {
		return "[debug] "
	}
	return ""
}

func init() {
	if DebugMode {
		SetLoggingLevel(Debug)
	} else {
		SetLoggingLevel(Info)
	}
}
