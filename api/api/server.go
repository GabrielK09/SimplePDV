package api

import (
	"log"
	"myApi/api/cors"
	authcontroller "myApi/api/handle/modules/controller/auth"
	cashregisterController "myApi/api/handle/modules/controller/cashRegister"
	customerController "myApi/api/handle/modules/controller/customer"
	dashboardcontroller "myApi/api/handle/modules/controller/dashBoard"
	payMentController "myApi/api/handle/modules/controller/payMentForms"
	productController "myApi/api/handle/modules/controller/product"
	reportController "myApi/api/handle/modules/controller/report"
	root "myApi/api/handle/modules/controller/root"
	saleController "myApi/api/handle/modules/controller/sale"
	u "myApi/helpers/logger"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer() {

	r := mux.NewRouter()

	cors := cors.WithCORS(r)

	r.HandleFunc("/api/root", root.HandleRoot)

	// Auth \\
	r.HandleFunc("/api/auth/login", authcontroller.HandleAuth).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Products \\
	r.HandleFunc("/api/products/all", productController.HandleGetProduct).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/products/find/{id}", productController.HandleGetByIdProduct).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/products/create", productController.HandlePostProduct).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/products/update/{id}", productController.HandlePutProduct).Methods(http.MethodPut, http.MethodOptions)
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

	// DashBoard \\
	r.HandleFunc("/api/dash-board/totales", dashboardcontroller.HandleProcessGetDashBoard).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/dash-board/popular-itens", dashboardcontroller.HandleProcessGetDashBoardPopularItens).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Relatórios \\
	r.HandleFunc("/api/report", reportController.HandlePostReports).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	log.Println("Servidor rodando em http://localhost:8000/api")
	u.ErrorLogger.Fatal(http.ListenAndServe(":8000", cors))
}
