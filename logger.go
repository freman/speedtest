package speedtest

import "fmt"

type Logger interface {
	Debugf(f string, v ...interface{})
	Warnf(f string, v ...interface{})
}

type voidLogger struct{}

func (l *voidLogger) Warnf(f string, v ...interface{})  {}
func (l *voidLogger) Debugf(f string, v ...interface{}) {}

type StdoutLogger struct{}

func (l *StdoutLogger) Warnf(f string, v ...interface{}) {
	fmt.Printf("[WARN]  "+f+"\n", v...)
}
func (l *StdoutLogger) Debugf(f string, v ...interface{}) {
	fmt.Printf("[DEBUG] "+f+"\n", v...)
}
