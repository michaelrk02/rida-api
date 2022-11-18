package api

import (
    "encoding/json"
    "errors"
    "log"
    "net/http"

    "github.com/michaelrk02/rida-api/service"
)

const (
    ErrUnauthorized string = "Unauthorized request"
    ErrForbidden string = "Access denied"
    ErrMalformedRequest string = "Malformed request"
    ErrInvalidRequest string = "Invalid request parameters"
    ErrNotFound string = "Resource not found"

    ErrServerSide string = "Server-side error"
)

var ErrorHandled error = errors.New("error handled")

type Error struct {
    Message string `json:"message"`
}

func (e Error) Send(w http.ResponseWriter, status int, source error) {
    app := service.GetApplication()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(&e)
    log.Printf("ERROR (%s) : %s\n", e.Message, source)

    if app.Params.Panic {
        panic(source)
    }
}
