package core

import (
	"fmt"
	"log"
)

type Logger interface {
	SetDebug(bool)
	Printf(string, ...any)
}

type DefaultPluginLogger struct {
	pluginName string
	debugMode  bool
}

func NewDefaultPluginLogger(pluginName string) *DefaultPluginLogger {
	return &DefaultPluginLogger{
		pluginName: pluginName,
		debugMode:  false,
	}
}

func (l *DefaultPluginLogger) SetDebug(v bool) {
	l.debugMode = v
}

func (l DefaultPluginLogger) Printf(format string, args ...any) {
	if l.debugMode {
		f := fmt.Sprintf("[PLUGIN: %s]: %s", l.pluginName, format)
		log.Printf(f, args...)
	}
}
