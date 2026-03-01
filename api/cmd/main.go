package main

import (
	"myApi/api"
	"myApi/db"
	loggerHelper "myApi/helpers/logger"
	u "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	"myApi/interface/dashBoard"
	paymentform "myApi/interface/payMentForm"
	"myApi/interface/pdv"
	"myApi/interface/product"
	"myApi/interface/sale"
)

func main() {
	loggerHelper.Logger()
	db, err := db.Init()

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao conectar ao banco: ", err)
	}

	u.GeneralLogger.Println("Banco de dados conectado com sucesso!")

	product.SetConnection(db)
	cashRegister.SetConnection(db)
	sale.SetConnection(db)
	customer.SetConnection(db)
	dashBoard.SetConnection(db)
	paymentform.SetConnection(db)
	pdv.SetConnection(db)

	if err = paymentform.CreateDefaultPayMents(); err != nil {
		u.ErrorLogger.Fatal("Erro ao criar as espécies padrão: ", err)
	}

	if err = customer.CreateDefaultCustomer(); err != nil {
		u.ErrorLogger.Fatal("Erro ao criar o consumidor padrão: ", err)
	}

	api.StartServer()
}
