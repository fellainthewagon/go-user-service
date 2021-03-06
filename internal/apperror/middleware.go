package apperror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var appError *AppError

		err := h(w, r)
		if err != nil {
			if errors.As(err, &appError) {
				switch {
				case errors.Is(err, NotFoundError):
					w.WriteHeader(http.StatusNotFound)
					w.Write(NotFoundError.MarshalError())
					return
				case errors.Is(err, BadRequestError):
					w.WriteHeader(http.StatusBadRequest)
					w.Write(BadRequestError.MarshalError())
					return
				default:
					appError = err.(*AppError)

					w.WriteHeader(http.StatusBadRequest)
					w.Write(appError.MarshalError())
					return
				}
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).MarshalError())
		}
	}
}
