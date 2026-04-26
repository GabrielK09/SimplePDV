package authcontroller

import (
	"context"
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/user"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")
var ctx = context.Background()

func checkPasswordHash(password, hash string) bool {
	u.InfoLogger.Printf("Comparando: %s - %s", password, hash)

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}

func createToken(userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": userName,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		u.ErrorLogger.Println("Erro ao gerar o token.", err)
		return "", err
	}

	return tokenString, nil
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)

		return
	}

	var bodyUserData user.UserContract

	if err := json.NewDecoder(r.Body).Decode(&bodyUserData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados:", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if err := bodyUserData.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := responsehelper.Response(false, err, "Campos obrigatórios ausentes.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	userByLogin, err := user.ShowByLogin(bodyUserData.Login)

	u.InfoLogger.Println("Valor de userByLogin", userByLogin)
	u.InfoLogger.Println("Valor de bodyUserData", bodyUserData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := responsehelper.Response(false, err, "Erro ao localizar o usuário.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if userByLogin == nil {
		w.WriteHeader(http.StatusUnauthorized)
		u.ErrorLogger.Println("Usuário não localizado")
		resp := responsehelper.Response(false, nil, "Credencias incorretas.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	isCorrectPassword := checkPasswordHash(
		bodyUserData.Password,
		userByLogin.Password,
	)

	if !isCorrectPassword {
		w.WriteHeader(http.StatusUnauthorized)
		u.ErrorLogger.Println("Senha incorreta.")
		resp := responsehelper.Response(false, nil, "Credencias incorretas.")

		json.NewEncoder(w).Encode(resp)
		return
	} else {
		tokenString, err := createToken(userByLogin.Name)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp := responsehelper.Response(false, nil, "Erro ao fazer login.")

			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := responsehelper.Response(false, tokenString, "Login bem sucedido!")

		json.NewEncoder(w).Encode(resp)
		return
	}
}
