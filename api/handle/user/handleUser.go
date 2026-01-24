package handle

import (
	"context"
	"encoding/json"
	"log"
	"myApi/db"
	responsehelper "myApi/helpers/response"
	usersContract "myApi/interface/users"
	"net/http"
	"time"
)

func HandleUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permitido.")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	if db.Pool == nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, nil, "Pool de conexões não inicializado.")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.Pool.Query(ctx, `
		SELECT 
			id, 
			name, 
			email 
		FROM 
			users 

		ORDER BY 
			id
		
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Falha ao consultar usuários:", err)
		resp := responsehelper.Response(false, nil, "Falha ao consultar usuários.")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	defer rows.Close()

	var users []usersContract.User

	for rows.Next() {
		var u usersContract.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Falha ao scanear usuário:", err)
			resp := responsehelper.Response(false, nil, "Erro ao ler dados do usuário.")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Erro durante iteração dos rows:", err)
		resp := responsehelper.Response(false, nil, "Erro durante a leitura dos dados.")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := responsehelper.Response(true, users, "Listando todos os usuários.")
	_ = json.NewEncoder(w).Encode(resp)
}
