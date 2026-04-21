package customercontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/customer"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func HandleGetCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	customers, err := customer.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao retornar todos os clientes."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, customers, "Todos os clientes.")

	json.NewEncoder(w).Encode(resp)
}

func HandleGetCustomerById(w http.ResponseWriter, r *http.Request) {
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

	customer, err := customer.Show(id)

	if customer == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Cliente não localizado."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, customer, "Cliente")

	json.NewEncoder(w).Encode(resp)
}

func HandlePostCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var c customer.CustomerContract

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao validar os dados.", err)
		resp := responsehelper.Response(false, err, "Erro ao validar os dados do cliente.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := c.Create(); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Erro ao cadastrar o cliente: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Erro ao cadastrar o cliente."),
		)

		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := responsehelper.Response(true, c, "Cliente cadastrado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}

func HandleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
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

	if id <= 1 {
		//O cliente padrão não pode ser desativado.
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O Id precisa ser maior que 1")
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if id == 1 {
		//O cliente padrão não pode ser desativado.
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O cliente padrão não pode ser desativado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, nil, "O cliente padrão não pode ser desativado."),
		)

		return
	}

	if err := customer.Delete(id, time.Now()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao deletar o cliente: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao deletar o Cliente."),
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Cliente deletado com sucesso!"))
}

func HandlePutCustomer(w http.ResponseWriter, r *http.Request) {
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

	if id <= 0 {
		//O cliente padrão não pode ser desativado.
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O Id precisa ser maior que 1")
		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Id inválido."))

		return
	}

	if id == 1 {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O cliente padrão não pode ser alterado.")

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "O cliente padrão não pode ser desativado."))

		return
	}

	u.SuccessLoger.Println("ID válidado")

	customer, err := customer.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Cliente não localizado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, nil, "Cliente não localizado."),
		)

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Erro ao processar dados:", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar dados"))

		return
	}

	customer.Id = id

	updated, err := customer.Update()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("O cliente padrão não pode ser alterado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	resp := responsehelper.Response(true, updated, "Cliente alterado com sucesso!")

	json.NewEncoder(w).Encode(resp)

}
