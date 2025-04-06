package api

import (
	"net/http"

	"github.com/bytedance/sonic"
)

type ErrorResponse struct {
	Status    int    `json:"cod"`
	ErrorDesc string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, status int, err error) {
	sonic.ConfigDefault.NewEncoder(w).Encode(ErrorResponse{
		Status:    status,
		ErrorDesc: err.Error(),
	})
}
