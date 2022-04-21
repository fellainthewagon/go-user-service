package user

import (
	"net/http"
	"rest-api/internal/handlers"

	"github.com/julienschmidt/httprouter"
)

// helper
// var _ handlers.IHandler = &handler{}

const (
	usersURL = "/users"
	userURL = "/users/:uuid"
)

type handler struct {
}

func NewHandler() handlers.IHandler {
	return &handler{}
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
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("Create user..."))
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
}

func (h *handler) PartUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("All users list..."))
}
