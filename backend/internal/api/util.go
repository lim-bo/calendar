package api

import (
	"log/slog"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Status    int    `json:"cod"`
	ErrorDesc string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, status int, err error) {
	jsonerr := sonic.ConfigDefault.NewEncoder(w).Encode(ErrorResponse{
		Status:    status,
		ErrorDesc: err.Error(),
	})
	if jsonerr != nil {
		slog.Error("sending error response issue", slog.String("error_value", jsonerr.Error()))
	}
}

func WriteLoginResponse(w http.ResponseWriter, uid uuid.UUID) {
	jsonerr := sonic.ConfigDefault.NewEncoder(w).Encode(map[string]interface{}{
		"cod": 200,
		"uid": uid.String(),
	})
	if jsonerr != nil {
		slog.Error("sending login response issue", slog.String("error_value", jsonerr.Error()))
	}
}
