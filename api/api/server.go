package api

import (
	"encoding/json"
	"fmt"
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
	responsehelper "myApi/helpers/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var secretKey = []byte("secret-key")

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Algortimo inválido: %s", t.Method.Alg())
		}

		return secretKey, nil
	})

	if err != nil {
		u.ErrorLogger.Println("Erro ao verificar o token.", err)
		return err
	}

	if !token.Valid {
		u.ErrorLogger.Println("Token inválido.")
		return fmt.Errorf("Token inválido.")
	}

	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			resp := responsehelper.Response(false, nil, "Token inválido.")

			json.NewEncoder(w).Encode(resp)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			resp := responsehelper.Response(false, nil, "Token inválido.")

			json.NewEncoder(w).Encode(resp)
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			resp := responsehelper.Response(false, nil, "Token inválido.")

			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := verifyToken(tokenString); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := responsehelper.Response(false, nil, "Token inválido.")

			json.NewEncoder(w).Encode(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StartServer() {
	r := mux.NewRouter()
	c := cors.WithCORS(r)

	r.HandleFunc("/api/ping", root.HandleRoot)
	r.HandleFunc("/api/auth/login", authcontroller.HandleAuth).Methods(http.MethodPost, http.MethodOptions)

	// Adiciona o prefixo, e usa a função handler como middleware
	proteced := r.PathPrefix("/api").Subrouter()
	proteced.Use(AuthMiddleware)

	// Products \\
	proteced.HandleFunc("/products/all", productController.HandleGetProduct).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/products/find/{id}", productController.HandleGetByIdProduct).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/products/find-by-name", productController.HandleGetByNameProduct).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/products/create", productController.HandlePostProduct).Methods(http.MethodPost, http.MethodOptions)
	proteced.HandleFunc("/products/update/{id}", productController.HandlePutProduct).Methods(http.MethodPut, http.MethodOptions)
	proteced.HandleFunc("/products/delete/{id}", productController.HandleDeleteProduct).Methods(http.MethodDelete, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Cash Register \\
	proteced.HandleFunc("/cash-register/all", cashregisterController.HandleGetCashRegister).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/cash-register/create", cashregisterController.HandlePostCashRegister).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Sale \\
	proteced.HandleFunc("/sale/all", saleController.HandleGetSale).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/sale/details/{id}", saleController.HandleGetSaleWithProducts).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/sale/cancel/{id}", saleController.HandlePutCancelSale).Methods(http.MethodPut, http.MethodOptions)
	proteced.HandleFunc("/sale/pay", saleController.HandlePutPaySale).Methods(http.MethodPut, http.MethodOptions)
	proteced.HandleFunc("/sale/create", saleController.HandlePostSale).Methods(http.MethodPost, http.MethodOptions)
	proteced.HandleFunc("/sale/new-itens", saleController.HandleInsertNewProducts).Methods(http.MethodPut, http.MethodOptions)
	proteced.HandleFunc("/sale/pay-ment-forms", payMentController.HandleGetPayMentForms).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/sale/update/pay-ment-forms/pix-key", payMentController.HandlePutPayMentForms).Methods(http.MethodPut, http.MethodOptions)

	// Customers \\
	proteced.HandleFunc("/customer/all", customerController.HandleGetCustomer).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/customer/find/{id}", customerController.HandleGetCustomerById).Methods(http.MethodGet, http.MethodOptions)
	proteced.HandleFunc("/customer/create", customerController.HandlePostCustomer).Methods(http.MethodPost, http.MethodOptions)
	proteced.HandleFunc("/customer/update/{id}", customerController.HandlePutCustomer).Methods(http.MethodPut, http.MethodOptions)
	proteced.HandleFunc("/customer/delete/{id}", customerController.HandleDeleteCustomer).Methods(http.MethodDelete, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// DashBoard \\
	proteced.HandleFunc("/dash-board/totales", dashboardcontroller.HandleProcessGetDashBoard).Methods(http.MethodPost, http.MethodOptions)
	proteced.HandleFunc("/dash-board/popular-itens", dashboardcontroller.HandleProcessGetDashBoardPopularItens).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Relatórios \\
	proteced.HandleFunc("/report", reportController.HandlePostReports).Methods(http.MethodPost, http.MethodOptions)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	log.Println("Servidor rodando em http://localhost:8000/api")
	u.ErrorLogger.Fatal(http.ListenAndServe(":8000", c))
}
