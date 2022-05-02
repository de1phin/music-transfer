package logger

import (
	"log"
	"os"
)

type FileLogger struct {
	file        *os.File
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	fileLogger := &FileLogger{}
	fileLogger.file = file
	fileLogger.errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	fileLogger.infoLogger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	return fileLogger, nil
}

func (fl *FileLogger) Error(x ...any) {
	fl.errorLogger.Println(x...)
}

func (fl *FileLogger) Info(x ...any) {
	fl.infoLogger.Println(x...)
}

func (fl *FileLogger) Close() error {
	return fl.file.Close()
}
