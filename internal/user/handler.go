package user

import (
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

// helper
// var _ handlers.IHandler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
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
	router.GET(usersURL, h.GetAllUsers)
	router.POST(usersURL, h.CreateUser)
	router.GET(userURL, h.GetUser)
	router.PUT(userURL, h.UpdateUser)
	router.PATCH(userURL, h.PartUpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("Create user..."))
	w.WriteHeader(201)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
}

func (h *handler) PartUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
	w.WriteHeader(200)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
	w.WriteHeader(204)
}
