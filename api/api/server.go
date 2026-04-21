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
	paymentcontroller "myApi/api/handle/modules/controller/payMent"
	paymentFormscontroller "myApi/api/handle/modules/controller/payMentForms"
	productController "myApi/api/handle/modules/controller/product"
	productcharacteristicscontroller "myApi/api/handle/modules/controller/product/productCharacteristics"
	reportController "myApi/api/handle/modules/controller/report"
	root "myApi/api/handle/modules/controller/root"
	saleController "myApi/api/handle/modules/controller/sale"
	shoppingcontroller "myApi/api/handle/modules/controller/shopping"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"net/http"
	"os"
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

/*
	func AuthMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
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
*/
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Token inválido."))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Token inválido."))
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Token inválido."))
			return
		}

		if err := verifyToken(tokenString); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Token inválido."))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StartServer() {
	r := mux.NewRouter()
	c := cors.WithCORS(r)

	api := r.PathPrefix("/api").Subrouter()

	r.HandleFunc("/api/ping", root.HandleRoot)
	r.HandleFunc("/api/auth/login", authcontroller.HandleAuth)

	api.Use(AuthMiddleware)

	// Products \\
	api.HandleFunc("/products/all", productController.HandleGetProduct)
	api.HandleFunc("/products/find/{id}", productController.HandleGetByIdProduct)
	api.HandleFunc("/products/find-by-name", productController.HandleGetByNameProduct)
	api.HandleFunc("/products/create", productController.HandlePostProduct)
	api.HandleFunc("/products/create/characteristics", productcharacteristicscontroller.HandlePostCreateProductCharacteristics)
	api.HandleFunc("/products/update/characteristics", productcharacteristicscontroller.HandlePutUpdateProductCharacteristics)
	api.HandleFunc("/products/{product_id}/find/characteristics/{grid_id}", productcharacteristicscontroller.HandleGetProductCharacteristicsByGridId)
	api.HandleFunc("/products/find/characteristics/{id}", productcharacteristicscontroller.HandleGetShowProductCharacteristics)
	api.HandleFunc("/products/{product_id}/delete/characteristics/{id}", productcharacteristicscontroller.HandleDeleteProductCharacteristics)
	api.HandleFunc("/products/update/{id}", productController.HandlePutProduct)
	api.HandleFunc("/products/delete/{id}", productController.HandleDeleteProduct)
	api.HandleFunc("/products/active/{id}", productController.HandleActiveProduct)
	api.HandleFunc("/products/verify-qtdes", productController.HandleVerifyQtdes)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Cash Register \\
	api.HandleFunc("/cash-register/all", cashregisterController.HandleGetCashRegister)
	api.HandleFunc("/cash-register/create", cashregisterController.HandlePostCashRegister)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Sale \\
	api.HandleFunc("/sale/create", saleController.HandlePostSale)
	api.HandleFunc("/sale/new-itens", saleController.HandleNewItens)
	api.HandleFunc("/sale/all", saleController.HandleGetSale)
	api.HandleFunc("/sale/details/{id}", saleController.HandleGetSaleWithProducts)

	// Customers \\
	api.HandleFunc("/customer/all", customerController.HandleGetCustomer)
	api.HandleFunc("/customer/find/{id}", customerController.HandleGetCustomerById)
	api.HandleFunc("/customer/create", customerController.HandlePostCustomer)
	api.HandleFunc("/customer/update/{id}", customerController.HandlePutCustomer)
	api.HandleFunc("/customer/delete/{id}", customerController.HandleDeleteCustomer)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// DashBoard \\
	api.HandleFunc("/dash-board/totales", dashboardcontroller.HandleProcessGetDashBoard)
	api.HandleFunc("/dash-board/popular-itens", dashboardcontroller.HandleProcessGetDashBoardPopularItens)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Relatórios \\
	api.HandleFunc("/report", reportController.HandlePostReports)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Compras \\
	api.HandleFunc("/shopping/all", shoppingcontroller.HandleGetAllShopping)
	api.HandleFunc("/shopping/return-last-load", shoppingcontroller.HandleGetLastShoppingLoad)
	api.HandleFunc("/shopping/create", shoppingcontroller.HandlePostCreateShopping)
	api.HandleFunc("/shopping/details/{id}", shoppingcontroller.HandleGetShoppingById)
	api.HandleFunc("/shopping/update", shoppingcontroller.HandlePutUpdateShopping)
	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	// Pagamentos \\
	api.HandleFunc("/pay-ment-forms/pay", paymentcontroller.HandlePutPaySaleOrShopping)

	api.HandleFunc("/pay-ment-forms/all", paymentFormscontroller.HandleGetPayMentForms)

	api.HandleFunc("/pay-ment-forms/update/pix-key", paymentFormscontroller.HandlePutPayMentForms)

	api.HandleFunc("/cancel-operation", paymentcontroller.HandlePutCancelSaleOrShopping)

	// -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == -- == \\

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	log.Println("Servidor rodando em http://localhost:8000/api")
	u.ErrorLogger.Fatal(http.ListenAndServe(":"+port, c))
}
