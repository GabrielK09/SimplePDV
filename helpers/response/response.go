package responsehelper

import (
	"myApi/interface/response"
)

func Response(status bool, data interface{}, message string) response.ResponseContract {
	return response.ResponseContract{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
