package logger

import (
	"log"
	"os"
)

var Debuglog *log.Logger

func init() {
	logFile, err := os.Create("./debug.log")
	// defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	Debuglog = log.New(logFile, "[Debug]", log.LstdFlags|log.Lshortfile)
}
