package handlers

import "github.com/julienschmidt/httprouter"

type IHandler interface {
	Register(router *httprouter.Router)
}
