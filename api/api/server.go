package api

import (
	"log"
	productController "myApi/api/handle/modules/controller/product"
	root "myApi/api/handle/modules/controller/root"
	cashregisterController "myApi/api/handle/modules/controller/sale/cashRegister"
	saleController "myApi/api/handle/modules/controller/sale/sale"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/root", root.HandleRoot)

	// Products \\
	r.HandleFunc("/api/products/all", productController.HandleGetProduct)
	r.HandleFunc("/api/products/create", productController.HandlePostProduct)
	r.HandleFunc("/api/products/update/{id}", productController.HandleUpdateProduct)
	r.HandleFunc("/api/products/delete/{id}", productController.HandleDeleteProduct)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Cash Register \\
	r.HandleFunc("/api/cash-register/all", cashregisterController.HandleGetCashRegister)
	r.HandleFunc("/api/cash-register/create", cashregisterController.HandlePostCashRegister)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Sale \\
	r.HandleFunc("/api/sale/create", saleController.HandlePostSale)
	r.HandleFunc("/api/sale/pay", saleController.HandlePutPaySale).Methods(http.MethodPut, http.MethodOptions)
	//r.HandleFunc("/api/sale/pay", )

	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	log.Println("Servidor rodando em http://localhost:8000/api")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
