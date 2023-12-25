package core

type Plugin interface {
	Setup(*Core) error
	Run(*Core) error
}

type PluginLoggable interface {
	Plugin
	SetLogger(Logger)
	SetDebug(bool)
}
