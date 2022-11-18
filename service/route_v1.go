package service

import (
    "net/http"
)

type RouteCollectionV1 interface {
    CreateTest(http.ResponseWriter, *http.Request)
    GetAllTest(http.ResponseWriter, *http.Request)
    GetTest(http.ResponseWriter, *http.Request)
    UpdateTest(http.ResponseWriter, *http.Request)
    DeleteTest(http.ResponseWriter, *http.Request)

    LoginAdmin(http.ResponseWriter, *http.Request)
    CreateAdmin(http.ResponseWriter, *http.Request)
    GetAllAdmin(http.ResponseWriter, *http.Request)
    GetAdmin(http.ResponseWriter, *http.Request)
    UpdateAdmin(http.ResponseWriter, *http.Request)
    DeleteAdmin(http.ResponseWriter, *http.Request)

    CreatePeneliti(http.ResponseWriter, *http.Request)
    GetAllPeneliti(http.ResponseWriter, *http.Request)
    GetPeneliti(http.ResponseWriter, *http.Request)
    UpdatePeneliti(http.ResponseWriter, *http.Request)
    DeletePeneliti(http.ResponseWriter, *http.Request)
}
