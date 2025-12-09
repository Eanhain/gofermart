package domain

type Logger interface {
	Warnln(args ...any)
	Infoln(args ...any)
}
