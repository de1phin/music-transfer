package logger

import (
	"log"
	"os"
)

type FileLogger struct {
	file          *os.File
	errorLogger   *log.Logger
	defaultLogger *log.Logger
}

func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	fileLogger := &FileLogger{}
	fileLogger.file = file
	fileLogger.errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	fileLogger.defaultLogger = log.New(file, "", log.Ldate|log.Ltime)
	return fileLogger, nil
}

func (fl *FileLogger) Log(x ...any) {
	hasError := false
	for _, i := range x {
		if _, ok := i.(error); ok {
			hasError = true
			break
		}
	}

	if hasError {
		fl.errorLogger.Println(x...)
	} else {
		fl.defaultLogger.Println(x...)
	}
}

func (fl *FileLogger) Close() error {
	return fl.file.Close()
}
