package main

import (
	"log"
	"myApi/api"
	"myApi/db"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal("Erro ao conectar ao banco: ", err)
	}

	log.Println("Banco de dados conectado com sucesso!")

	api.StartServer()
}
