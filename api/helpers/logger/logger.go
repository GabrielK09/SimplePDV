package loggerHelper

import (
	"log"
	"os"
	"path/filepath"
)

var GeneralLogger *log.Logger
var InfoLogger *log.Logger
var ErrorLogger *log.Logger
var SuccessLoger *log.Logger

func Logger() {
	logDir := "log"

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatal("Erro ao criar o direiatório de log: ", err)
	}

	logPath := filepath.Join(logDir, "general-log.log")

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Erro ao abrir a log:", err)
	}

	GeneralLogger = log.New(file, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(file, "Info Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	SuccessLoger = log.New(file, "Success Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
