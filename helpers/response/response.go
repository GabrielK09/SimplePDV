package responsehelper

import (
	"myApi/interface/response"
)

func Response(status bool, data any, message string) response.ResponseContract {
	return response.ResponseContract{
		Status:  status,
		Data:    data,
		Message: message,
	}
}
