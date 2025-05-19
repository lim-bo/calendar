package api

import (
	"log/slog"
	"net/http"
	"net/mail"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/lim-bo/calendar/backend/models"
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

func WriteGetProfileResponse(w http.ResponseWriter, info *models.UserCredentialsRegister) {
	err := sonic.ConfigDefault.NewEncoder(w).Encode(info)
	if err != nil {
		slog.Error("sending profile info error", slog.String("error_value", err.Error()))
	}
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateEmailMult(emails []string) bool {
	for _, e := range emails {
		if !ValidateEmail(e) {
			return false
		}
	}
	return true
}
