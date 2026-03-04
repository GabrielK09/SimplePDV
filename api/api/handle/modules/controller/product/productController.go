package productcontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/product"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
)

func HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	products, err := product.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todos os produtos.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, products, "Todos os produtos cadastrados.")

	json.NewEncoder(w).Encode(resp)
}

func HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)

		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if err := product.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao deletar o produto."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, nil, "Produto deletado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}

func HandleGetByIdProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	products, err := product.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar o produtos.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, products, "Produto")

	json.NewEncoder(w).Encode(resp)
}

func HandlePostProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var payload product.ProductContract

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payload.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := responsehelper.Response(false, err, "Campos obrigatórios ausentes.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payload.Create(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao gravar o produto.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := responsehelper.Response(true, payload, "Produto criado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}

func HandlePutProduct(w http.ResponseWriter, r *http.Request) {
	var productData *product.ProductContract

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	productData, err = product.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		u.ErrorLogger.Println("Produto não localizado:", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Produto não localizado."),
		)

		return
	}

	u.InfoLogger.Println("Produto:", productData)

	if err := json.NewDecoder(r.Body).Decode(&productData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Erro ao processar dados: ", err)

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	productData.Id = id

	updated, err := productData.Update()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("Erro ao alterar o produto: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao atualizar"),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		responsehelper.Response(true, updated, "Produto alterado com sucesso!"),
	)
}
