package log

type Logger interface {
	Error(...any)
	Info(...any)
}
