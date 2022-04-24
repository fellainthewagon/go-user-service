package user

import (
	"net/http"
	"rest-api/internal/apperror"
	"rest-api/internal/handlers"
	"rest-api/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/users"
	userURL  = "/users/:id"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.IHandler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetAllUsers))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUser))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
	return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Create user..."))
	w.WriteHeader(201)
	return nil
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("All users list..."))
	w.WriteHeader(204)
	return nil
}
