package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Middleware func(http.ResponseWriter, *http.Request) error

	Error struct {
		StatusCode int    `json:"-"`
		Message    string `json:"message"`
		Details    string `json:"details"`
	}
)

func (e *Error) Error() string {
	return fmt.Sprintf(e.Message, ":", e.Details)
}

func (fn Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		errResp, ok := err.(*Error)
		if !ok {
			errResp = &Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "unexpected internal server error",
				Details:    err.Error(),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errResp.StatusCode)

		body, _ := json.Marshal(errResp)
		w.Write(body)
	}
}
