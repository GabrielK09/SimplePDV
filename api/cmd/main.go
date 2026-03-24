package main

import (
	"context"
	"flag"
	"log"
	"myApi/api"
	cashregisterController "myApi/api/handle/modules/controller/cashRegister"
	processpayment "myApi/api/services/payMent"
	reports "myApi/api/services/reports"
	"myApi/db"
	dbaction "myApi/db/dbActions/seeder"
	u "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	"myApi/interface/dashBoard"
	paymentform "myApi/interface/payMentForm"
	"myApi/interface/pdv"
	"myApi/interface/product"
	productcharacteristics "myApi/interface/product/productCharacteristics"
	"myApi/interface/sale"
	"myApi/interface/shopping"
	"myApi/interface/user"
	"myApi/jobs"
	"os"
)

var ctx = context.Background()
var jobFlag string
var dbFlag string

func main() {
	wd, _ := os.Getwd()

	log.Println("Rodando em: ", wd)

	u.Logger()
	dbConn, err := db.Init()

	if err != nil {
		log.Fatal("Erro ao conectar ao banco: ", err)
	}

	flag.StringVar(&dbFlag, "db", "", "confirm if is a db flag")
	flag.StringVar(&jobFlag, "job", "", "confirm if is a job")

	flag.Parse()

	switch jobFlag {
	case "createUser":
		jobs.CreateUser(dbConn, ctx)
		return

	case "resetSite":
		jobs.ResetSite(dbConn, ctx)
		jobs.CreateUser(dbConn, ctx)
		return
	}

	switch dbFlag {
	case "seed":
		dbaction.DBSeed(dbConn, ctx)

		return
	}

	log.Println("Banco de dados conectado com sucesso!")

	product.SetConnection(dbConn)
	productcharacteristics.SetConnection(dbConn)
	cashregisterController.SetConnection(dbConn) // For manual insert
	cashRegister.SetConnection(dbConn)
	sale.SetConnection(dbConn)
	customer.SetConnection(dbConn)
	dashBoard.SetConnection(dbConn)
	paymentform.SetConnection(dbConn)
	pdv.SetConnection(dbConn)
	reports.SetConnection(dbConn)
	shopping.SetConnection(dbConn)
	user.SetConnection(dbConn)
	processpayment.SetConnection(dbConn)

	if err = paymentform.CreateDefaultPayMents(); err != nil {
		log.Fatal("Erro ao criar as espécies padrão: ", err)
	}

	if err = customer.CreateDefaultCustomer(); err != nil {
		log.Fatal("Erro ao criar o consumidor padrão: ", err)
	}

	api.StartServer()
}
