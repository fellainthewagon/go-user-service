package user

import (
	"context"
	"encoding/json"
	"net/http"
	"rest-api/internal/apperror"
	"rest-api/internal/handlers"
	"rest-api/pkg/logging"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/users"
	userURL  = "/users/:id"
)

type handler struct {
	service Service
	logger  *logging.Logger
}

func NewHandler(service Service, logger *logging.Logger) handlers.IHandler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetAllUsers))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUser))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var dto CreateUserDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return apperror.BadRequestError
	}

	user, err := h.service.Create(context.TODO(), dto)
	if err != nil {
		return err
	}

	userBytes, err := user.Marshal()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userBytes)
	return nil
}

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.GetAllUsers(context.TODO())
	if err != nil {
		return err
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersBytes)

	return nil
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id, err := h.getUserId(r)
	if err != nil {
		return err
	}

	user, err := h.service.GetUser(context.TODO(), id)
	if err != nil {
		return err
	}

	userBytes, err := user.Marshal()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userBytes))

	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id, err := h.getUserId(r)
	if err != nil {
		return err
	}

	var dto UpdateUserDTO
	dto.ID = id

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return apperror.BadRequestError
	}

	err = h.service.UpdateUser(context.TODO(), dto)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := h.getUserId(r)
	if err != nil {
		return err
	}

	err = h.service.DeleteUser(context.TODO(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
	return nil
}

func (h *handler) getUserId(r *http.Request) (string, error) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	if id == "" {
		return "", apperror.BadRequestError
	}

	return id, nil
}