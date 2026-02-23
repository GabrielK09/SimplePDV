package main

import (
	"myApi/api"
	"myApi/db"
	loggerHelper "myApi/helpers/logger"
	u "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	paymentform "myApi/interface/payMentForm"
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
	paymentform.SetConnection(db)
	customer.SetConnection(db)

	api.StartServer()
}
