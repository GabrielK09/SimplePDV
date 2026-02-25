package api

import (
	"log"
	"myApi/api/cors"
	customerController "myApi/api/handle/modules/controller/customer"
	productController "myApi/api/handle/modules/controller/product"
	root "myApi/api/handle/modules/controller/root"
	cashregisterController "myApi/api/handle/modules/controller/sale/cashRegister"
	payMentController "myApi/api/handle/modules/controller/sale/payMentForms"
	saleController "myApi/api/handle/modules/controller/sale/sale"
	u "myApi/helpers/logger"
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
	r.HandleFunc("/api/sale/all", saleController.HandleGetSale).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/sale/details/{id}", saleController.HandleGetSaleWithProducts).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/sale/cancel/{id}", saleController.HandlePutCancelSale).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/sale/pay", saleController.HandlePutPaySale).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/sale/create", saleController.HandlePostSale).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/sale/pay-ment-forms", payMentController.HandleGetPayMentForms).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/sale/update/pay-ment-forms/pix-key", payMentController.HandlePutPayMentForms).Methods(http.MethodPut, http.MethodOptions)

	// Customers \\
	r.HandleFunc("/api/customer/all", customerController.HandleGetCustomer).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/customer/find/{id}", customerController.HandleGetCustomerById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/customer/create", customerController.HandlePostCustomer).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/customer/update/{id}", customerController.HandlePutCustomer).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/customer/delete/{id}", customerController.HandleDeleteCustomer).Methods(http.MethodDelete, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	log.Println("Servidor rodando em http://localhost:8000/api")
	u.ErrorLogger.Fatal(http.ListenAndServe(":8000", cors))
}
