package main

import "log"

type Logger interface {
	LogInfo(message string)
	LogError(message string)
}

type ConsoleLogger struct {
	enabled bool
}

func NewConsoleLogger(enabled bool) *ConsoleLogger {
	return &ConsoleLogger{enabled: enabled}
}

func (l *ConsoleLogger) LogInfo(message string) {
	if l.enabled {
		log.Printf("INFO: %s\n", message)
	}
}

func (l *ConsoleLogger) LogError(message string) {
	if l.enabled {
		log.Printf("ERROR: %s\n", message)
	}
}
