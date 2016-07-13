package utils

type Logger interface {
	Printf(format string, args ...interface{})
}
