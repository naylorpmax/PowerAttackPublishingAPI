package middleware

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
	fmt.Println("new request")
	err := fn(w, r)
	if err != nil {
		errResp, ok := err.(*Error)
		if !ok {
			// TODO: add default error response
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errResp.StatusCode)

		body, err := json.Marshal(errResp)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(body)
	}

	fmt.Println("request handled")
}
