package loggerHelper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var GeneralLogger *log.Logger
var InfoLogger *log.Logger
var ErrorLogger *log.Logger

func Logger() {
	absPath, err := filepath.Abs("../api/log")

	if err != nil {
		fmt.Println("Erro ao pegar o caminho: ", err)
	}

	generalLog, err := os.OpenFile(absPath+"/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Erro ao abrir:", err)
		os.Exit(1)
	}

	errorLog, err := os.OpenFile(absPath+"/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Erro ao abrir:", err)
		os.Exit(1)
	}

	infoLog, err := os.OpenFile(absPath+"/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Erro ao abrir:", err)
		os.Exit(1)
	}

	GeneralLogger = log.New(generalLog, "General Logger: \t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorLog, "Error Logger: \t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(infoLog, "Info Logger: \t", log.Ldate|log.Ltime|log.Lshortfile)
}
