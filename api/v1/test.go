package v1

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/michaelrk02/rida-api/api"
    "github.com/michaelrk02/rida-api/model"
    "github.com/michaelrk02/rida-api/resource"
    "gorm.io/gorm"
)

func (routes *RouteCollection) CreateTest(w http.ResponseWriter, r *http.Request) {
    var err error

    var req resource.TestRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    test := new(model.Test).FromRequest(&req)
    err = routes.App.DB.Create(test).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    resp := test.ToResponseOnCreate()
    w.WriteHeader(201)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) GetAllTest(w http.ResponseWriter, r *http.Request) {
    var err error

    var tests []model.Test
    err = routes.App.DB.Find(&tests).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    var resp resource.TestResponseCollection
    resp.Count = len(tests)
    resp.Data = make([]resource.TestResponse, resp.Count)
    for i := range tests {
        resp.Data[i] = tests[i].ToResponse()
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) GetTest(w http.ResponseWriter, r *http.Request) {
    var err error

    testID := chi.URLParam(r, "test")

    var test model.Test
    err = routes.App.DB.Where("id = ?", testID).First(&test).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
        } else {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        }
        return
    }

    resp := test.ToResponse()
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) UpdateTest(w http.ResponseWriter, r *http.Request) {
    var err error

    testID := chi.URLParam(r, "test")

    var req resource.TestRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    test := new(model.Test).FromRequest(&req).AssignID(testID)
    err = routes.App.DB.Save(test).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    w.WriteHeader(200)
}

func (routes *RouteCollection) DeleteTest(w http.ResponseWriter, r *http.Request) {
    var err error

    testID := chi.URLParam(r, "test")

    err = routes.App.DB.Where("id = ?", testID).Delete(&model.Test{}).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    w.WriteHeader(200)
}
