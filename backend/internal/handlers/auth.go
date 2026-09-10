package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"example.com/german/backend/internal/service"
)

const maxLoginRequestBytes = 16 * 1024

type authService interface {
	Login(context.Context, string, string) (service.LoginResponse, error)
	Authenticate(context.Context, string) (service.AuthUser, error)
}

type AuthHandler struct {
	service authService
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthHandler(authService authService) (*AuthHandler, error) {
	if authService == nil {
		return nil, errors.New("auth service is required")
	}
	return &AuthHandler{service: authService}, nil
}

func (handler *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", handler.login)
	mux.HandleFunc("POST /api/auth/logout", handler.logout)
	mux.HandleFunc("GET /api/auth/me", handler.me)
}

func (handler *AuthHandler) login(w http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(w, request.Body, maxLoginRequestBytes)
	var payload loginRequest
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		writeJSONRequestError(w, err)
		return
	}
	if err := ensureJSONEnd(decoder); err != nil {
		writeJSONRequestError(w, err)
		return
	}
	response, err := handler.service.Login(request.Context(), payload.Username, payload.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeAPIError(w, http.StatusUnauthorized, "invalid_credentials", "invalid username or password")
			return
		}
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *AuthHandler) logout(w http.ResponseWriter, request *http.Request) {
	if _, err := handler.authenticate(request); err != nil {
		writeAuthError(w, err)
		return
	}
	// JWTs are stateless. Logout means the client discards its token.
	w.WriteHeader(http.StatusNoContent)
}

func (handler *AuthHandler) me(w http.ResponseWriter, request *http.Request) {
	user, err := handler.authenticate(request)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]service.AuthUser{"user": user})
}

func (handler *AuthHandler) authenticate(request *http.Request) (service.AuthUser, error) {
	return handler.service.Authenticate(request.Context(), request.Header.Get("Authorization"))
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidToken):
		w.Header().Set("WWW-Authenticate", "Bearer")
		writeAPIError(w, http.StatusUnauthorized, "unauthorized", "valid bearer token is required")
	case errors.Is(err, service.ErrForbidden):
		writeAPIError(w, http.StatusForbidden, "forbidden", "operation is forbidden")
	default:
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
	}
}

func decodeJSONBody(w http.ResponseWriter, request *http.Request, limit int64, destination any) error {
	request.Body = http.MaxBytesReader(w, request.Body, limit)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is required")
		}
		return err
	}
	return ensureJSONEnd(decoder)
}
