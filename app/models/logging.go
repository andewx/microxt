package models

import (
	"log"
)

// Logging framework
type Logging struct {
	errLog   *log.Logger
	warnLog  *log.Logger
	infoLog  *log.Logger
	debugLog *log.Logger
}

func NewLogging() *Logging {
	return &Logging{errLog: log.Default(), warnLog: log.Default(), infoLog: log.Default(), debugLog: log.Default()}
}
