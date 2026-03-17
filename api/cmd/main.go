package main

import (
	"context"
	"flag"
	"myApi/api"
	cashregisterController "myApi/api/handle/modules/controller/cashRegister"
	processpayment "myApi/api/services/payMent"
	reports "myApi/api/services/reports"
	"myApi/db"
	dbaction "myApi/db/dbActions/seeder"
	loggerHelper "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	"myApi/interface/dashBoard"
	paymentform "myApi/interface/payMentForm"
	"myApi/interface/pdv"
	"myApi/interface/product"
	"myApi/interface/sale"
	"myApi/interface/shopping"
	"myApi/interface/user"
	"myApi/jobs"
)

var ctx = context.Background()
var jobFlag string
var dbFlag string

func main() {
	loggerHelper.Logger()
	db, err := db.Init()

	if err != nil {
		loggerHelper.ErrorLogger.Fatal("Erro ao conectar ao banco: ", err)
	}

	flag.StringVar(&dbFlag, "db", "", "confirm if is a db flag")
	flag.StringVar(&jobFlag, "job", "", "confirm if is a job")

	flag.Parse()

	switch jobFlag {
	case "createUser":
		jobs.CreateUser(db, ctx)
		return

	case "resetSite":
		jobs.ResetSite(db, ctx)
		return
	}

	switch dbFlag {
	case "seed":
		dbaction.DBSeed(db, ctx)
		return
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
	shopping.SetConnection(db)
	user.SetConnection(db)
	processpayment.SetConnection(db)

	if err = paymentform.CreateDefaultPayMents(); err != nil {
		loggerHelper.ErrorLogger.Fatal("Erro ao criar as espécies padrão: ", err)
	}

	if err = customer.CreateDefaultCustomer(); err != nil {
		loggerHelper.ErrorLogger.Fatal("Erro ao criar o consumidor padrão: ", err)
	}

	api.StartServer()
}
