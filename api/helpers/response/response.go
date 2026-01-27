package responsehelper

import (
	"log"
	"myApi/interface/response"
)

func Response(status bool, data any, message string) response.ResponseContract {
	log.Println("Data recebido:", data)
	return response.ResponseContract{
		Status:  status,
		Data:    data,
		Message: message,
	}
}
