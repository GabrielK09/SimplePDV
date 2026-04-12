package shoppingcontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	productcharacteristics "myApi/interface/product/productCharacteristics"
	"myApi/interface/shopping"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductWithCharacteristics struct {
	Product         shopping.ShoppingItenContract                           `json:"product"`
	Characteristics []productcharacteristics.ProductCharacteristicsContract `json:"product_with_characteristics"`
}

func HandleGetAllShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	shoppings, err := shopping.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todos os produtos.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, shoppings, "Todas as compras cadastradas."))
}

func HandleGetLastShoppingLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	shoppingId, err := shopping.ReturnLastShoppingLoad()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao retornar o ID: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Falha na operação."))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, shoppingId, ""))
}

func HandlePostCreateShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	} // Erro de método da rota

	var payload shopping.ShoppingContract

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados: ", err)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payload.Validate(); len(err) > 0 {
		if err["load"] != "" {
			w.WriteHeader(http.StatusInternalServerError)
			u.ErrorLogger.Println("Carga duplicada.", err)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Dados duplicados"))
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Campos obrigatórios ausentes.", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Campos obrigatórios ausentes."))
		return
	}

	shoppingId, err := payload.Create()

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Erro ao cadastrar a compra.", err)
		resp := responsehelper.Response(false, err, "Erro ao cadastrar a compra.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responsehelper.Response(true, shoppingId, "Compra cadastrada com sucesso!"))
}

func HandlePutUpdateShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	} // Erro de método da rota

	var shoppingDetails shopping.ShoppingContract

	if err := json.NewDecoder(r.Body).Decode(&shoppingDetails); err != nil {
		u.ErrorLogger.Println("Erro ao ler os dados da edição da compra:", err)
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao ler os dados da edição da compra."))
		return
	}

	if shoppingDetails.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "ID da compra inválido."))
		return
	}

	if err := shoppingDetails.UpdateShopping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao alterar os dados da compra."))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Compra alterada com sucesso!"))
}

func HandleGetShoppingById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	} // Erro de método da rota

	params := mux.Vars(r)
	shoppingId, err := strconv.Atoi(params["id"])

	shoppingData, err := shopping.Show(shoppingId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("Erro ao buscar os dados da compra: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao buscar os dados da compra."))
		return
	}

	shoppingIntesData, err := shopping.ShowShoppingItens(shoppingId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("Erro ao buscar os itens da compra: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao buscar os itens da compra."))
		return
	}

	productCharacteristics, err := shopping.ShowShoppingGridItens(shoppingId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar todos os produtos."))
		return
	}

	var productsWithCharacteristics []ProductWithCharacteristics

	for _, p := range *shoppingIntesData {
		productsWithCharacteristics = append(productsWithCharacteristics, ProductWithCharacteristics{
			Product:         p,
			Characteristics: *productCharacteristics,
		})
	}

	repsonse := map[string]any{
		"shopping":             shoppingData,
		"shoppingWithProducts": productsWithCharacteristics,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, repsonse, "Forma de pagamento alterada com sucesso!"))
}
