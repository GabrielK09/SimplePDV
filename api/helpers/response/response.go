package responsehelper

import (
	u "myApi/helpers/logger"
	"myApi/interface/response"
)

func Response(status bool, data any, message string) response.ResponseContract {
	u.GeneralLogger.Println("Data recebido:", data)

	return response.ResponseContract{
		Status:  status,
		Data:    data,
		Message: message,
	}
}
