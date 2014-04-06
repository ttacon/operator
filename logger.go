package operator

import "log"

const operatorPrefix = "[operator] "

var logger operatorLogger

type operatorLogger struct{}

func (l operatorLogger) Log(args ...interface{}) {
	// TODO(ttacon): this is an ugly hack, make this better...
	args = append([]interface{}{operatorPrefix}, args...)
	log.Print(args...)
}

func (l operatorLogger) Logf(format string, args ...interface{}) {
	log.Printf(operatorPrefix+format, args...)
}
