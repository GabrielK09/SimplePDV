package api

import (
	"log"
	product "myApi/api/handle/modules/product"
	root "myApi/api/handle/modules/root"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/api", root.HandleRoot)
	http.HandleFunc("/api/product/register", product.HandleRegisterProduct)
	http.HandleFunc("/api/product/all", product.HandleGetAllProduct)
	http.HandleFunc("/api/product/update/id", product.HandleUpdateProduct)
	http.HandleFunc("/api/product/delete/id", product.HandleDeleteProduct)

	log.Println("Servidor rodando em http://localhost:9000/api")

	log.Fatal(http.ListenAndServe(":9000", nil))

}
