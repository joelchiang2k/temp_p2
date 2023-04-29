package logger

import (
	"log"
	"os"
	"sync"
)

type _logger struct{
	fileName string
	*log.Logger
}

var logger *_logger
var once sync.Once

func GetInst() *_logger {
	once.Do(func() {
		logger = createLog("LogFile.txt")
	})
	return logger
}

func createLog(fileName string) *_logger {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return &_logger{
		fileName: fileName,
		Logger: log.New(file, "Proj2", log.Lshortfile),
	}
}