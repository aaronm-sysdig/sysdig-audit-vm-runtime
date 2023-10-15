package DebugLog

import (
	"log"
)

var Debug bool

type DebugLogger struct{}

func (d *DebugLogger) Println(args ...interface{}) {
	if Debug {
		log.Println(args...)
	}
}

func (d *DebugLogger) Printf(format string, args ...interface{}) {
	if Debug {
		log.Printf(format, args...)
	}
}

func (d *DebugLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
