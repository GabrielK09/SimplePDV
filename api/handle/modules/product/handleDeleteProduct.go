package handle

import (
	"encoding/json"
	"log"
	responsehelper "myApi/helpers/response"
	"myApi/interface/product"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
)

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

		log.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if err := product.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Erro ao deletar o produto."),
		)

		return
	}

	resp := responsehelper.Response(true, nil, "Produto deletado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
