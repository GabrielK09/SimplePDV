package api

import (
	"log"
	product "myApi/api/handle/modules/product"
	root "myApi/api/handle/modules/root"
	cashregister "myApi/api/handle/modules/sale/cashRegister"
	sale "myApi/api/handle/modules/sale/sale"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api", root.HandleRoot)

	// Products \\
	r.HandleFunc("/api/product/all", product.HandleGetProduct)
	r.HandleFunc("/api/product/register", product.HandlePostProduct)
	r.HandleFunc("/api/product/update/{id}", product.HandleUpdateProduct)
	r.HandleFunc("/api/product/delete/{id}", product.HandleDeleteProduct)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Cash Register \\
	r.HandleFunc("/api/cash-register/all", cashregister.HandleGetCashRegister)
	r.HandleFunc("/api/cash-register/register", cashregister.HandlePostCashRegister)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Sale \\
	r.HandleFunc("/api/sale/register", sale.HandlePostSale)
	//r.HandleFunc("/api/sale/pay", )

	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	log.Println("Servidor rodando em http://localhost:9000/api")
	log.Fatal(http.ListenAndServe(":9000", r))
}
