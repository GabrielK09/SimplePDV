package api

import (
	"log"
	auth "myApi/api/handle/auth"
	root "myApi/api/handle/root"
	user "myApi/api/handle/user"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/api", root.HandleRoot)
	http.HandleFunc("/api/user/all", user.HandleUser)
	http.HandleFunc("/api/register", auth.HandleAuth)

	log.Println("Servidor rodando em http://localhost:9000/api")

	log.Fatal(http.ListenAndServe(":9000", nil))

}
