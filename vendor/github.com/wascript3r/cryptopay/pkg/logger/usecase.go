package logger

type Usecase interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
