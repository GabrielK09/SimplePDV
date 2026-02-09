package api

import (
	"log"
	"myApi/api/cors"
	productController "myApi/api/handle/modules/controller/product"
	root "myApi/api/handle/modules/controller/root"
	cashregisterController "myApi/api/handle/modules/controller/sale/cashRegister"
	saleController "myApi/api/handle/modules/controller/sale/sale"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {

	r := mux.NewRouter()

	cors := cors.WithCORS(r)

	r.HandleFunc("/api/root", root.HandleRoot)

	// Products \\
	r.HandleFunc("/api/products/all", productController.HandleGetProduct).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/products/find/{id}", productController.HandleGetByIdProduct).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/products/create", productController.HandlePostProduct).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/products/update/{id}", productController.HandleUpdateProduct).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/products/delete/{id}", productController.HandleDeleteProduct).Methods(http.MethodDelete, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Cash Register \\
	r.HandleFunc("/api/cash-register/all", cashregisterController.HandleGetCashRegister).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/cash-register/create", cashregisterController.HandlePostCashRegister).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Sale \\
	r.HandleFunc("/api/sale/pay", saleController.HandlePutPaySale).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/sale/create", saleController.HandlePostSale).Methods(http.MethodPost, http.MethodOptions)
	//r.HandleFunc("/api/sale/pay", )

	// Shopping \\

	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	log.Println("Servidor rodando em http://localhost:8000/api")
	log.Fatal(http.ListenAndServe(":8000", cors))
}
