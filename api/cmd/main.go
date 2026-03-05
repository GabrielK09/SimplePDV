package main

import (
	"flag"
	"myApi/api"
	cashregisterController "myApi/api/handle/modules/controller/cashRegister"
	"myApi/db"
	loggerHelper "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	"myApi/interface/dashBoard"
	paymentform "myApi/interface/payMentForm"
	"myApi/interface/pdv"
	"myApi/interface/product"
	"myApi/interface/reports"
	"myApi/interface/sale"
	"myApi/interface/user"
	"myApi/jobs"
	"os"
)

var job bool

func main() {
	loggerHelper.Logger()
	db, err := db.Init()

	if err != nil {
		loggerHelper.ErrorLogger.Fatal("Erro ao conectar ao banco: ", err)
	}

	flag.BoolVar(&job, "isJob", false, "confire if is a job")
	flag.Parse()

	if job {
		loggerHelper.InfoLogger.Println("É um job.")
		jobs.CreateUser()

		os.Exit(0)
	}

	loggerHelper.GeneralLogger.Println("Banco de dados conectado com sucesso!")

	product.SetConnection(db)
	cashregisterController.SetConnection(db) // For manual insert
	cashRegister.SetConnection(db)
	sale.SetConnection(db)
	customer.SetConnection(db)
	dashBoard.SetConnection(db)
	paymentform.SetConnection(db)
	pdv.SetConnection(db)
	reports.SetConnection(db)
	user.SetConnection(db)

	if err = paymentform.CreateDefaultPayMents(); err != nil {
		loggerHelper.ErrorLogger.Fatal("Erro ao criar as espécies padrão: ", err)
	}

	if err = customer.CreateDefaultCustomer(); err != nil {
		loggerHelper.ErrorLogger.Fatal("Erro ao criar o consumidor padrão: ", err)
	}

	api.StartServer()
}
