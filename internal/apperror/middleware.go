package apperror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appError *AppError

		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			if errors.As(err, &appError) {
				switch {
				case errors.Is(err, NotFoundError):
					w.WriteHeader(http.StatusNotFound)
					w.Write(NotFoundError.MarshalError())
					break
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
