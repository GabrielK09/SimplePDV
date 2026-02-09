package main

import (
	"log"
	"myApi/api"
	"myApi/db"
	"myApi/interface/cashRegister"
	paymentform "myApi/interface/payMentForm"
	"myApi/interface/product"
	"myApi/interface/sale"
)

func main() {
	db, err := db.Init()

	if err != nil {
		log.Fatal("Erro ao conectar ao banco: ", err)
	}

	log.Println("Banco de dados conectado com sucesso!")

	product.SetConnection(db)
	cashRegister.SetConnection(db)
	sale.SetConnection(db)
	paymentform.SetConnection(db)

	api.StartServer()
}
