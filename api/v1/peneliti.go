package v1

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/michaelrk02/rida-api/api"
    "github.com/michaelrk02/rida-api/model"
    "github.com/michaelrk02/rida-api/resource"
    "github.com/michaelrk02/rida-api/service"
    "gorm.io/gorm"
)

func (routes *RouteCollection) CreatePeneliti(w http.ResponseWriter, r *http.Request) {
    var err error

    var req resource.PenelitiRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    var resp resource.PenelitiResponseOnCreate

    err = routes.App.DB.Transaction(func(tx *gorm.DB) error {
        var err error

        var peneliti model.Peneliti

        err = tx.Where("nidn = ?", req.Nidn).First(&peneliti).Error
        if err == nil {
            api.Error{Message: "Peneliti already exists"}.Send(w, 400, err)
            return api.ErrorHandled
        }
        if err != nil {
            if !errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
                return api.ErrorHandled
            }
        }

        admin := r.Context().Value("auth_entity").(model.Admin)

        peneliti = model.Peneliti{}
        peneliti.FromRequest(&req)
        peneliti.DiciptakanOlehID = admin.ID

        if !admin.FakultasID.Valid {
            if req.FakultasID == nil {
                api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, errors.New("Fakultas must be present"))
                return api.ErrorHandled
            }
            peneliti.FakultasID = *req.FakultasID
        } else {
            peneliti.FakultasID = admin.FakultasID.String
        }

        err = routes.App.DB.Create(&peneliti).Error
        if err != nil {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            return api.ErrorHandled
        }

        resp = peneliti.ToResponseOnCreate()

        return nil
    })
    if err != nil {
        if !errors.Is(err, api.ErrorHandled) {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        }
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) GetAllPeneliti(w http.ResponseWriter, r *http.Request) {
    var err error

    ds := service.NewDataSource(
        100,
        []string{"peneliti.nidn", "peneliti.nama", "peneliti.scopus_author_id", "peneliti.gscholar_author_id", "Fakultas.nama", "DiciptakanOleh.nama"},
        []string{"peneliti.nidn", "peneliti.nama", "Fakultas.nama", "DiciptakanOleh.nama"},
        "peneliti.nama",
    ).FromRequest(r, "")

    var totalItems int64
    err = routes.App.DB.Model(&model.Peneliti{}).Joins("Fakultas").Joins("DiciptakanOleh").Scopes(ds.EnumerationScope).Count(&totalItems).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    ds.Populate(int(totalItems))

    err = ds.Validate()
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    var penelitiList []model.Peneliti
    err = routes.App.DB.Joins("Fakultas").Joins("DiciptakanOleh").Scopes(ds.PopulationScope).Find(&penelitiList).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    var resp resource.PenelitiResponseCollection
    resp.Population = ds.TotalItems
    resp.Display = ds.Display
    resp.Page = ds.Page
    resp.MaxPage = ds.MaxPage
    resp.Data = make([]resource.PenelitiResponse, len(penelitiList))
    for i := range penelitiList {
        resp.Data[i] = penelitiList[i].ToResponse()
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) GetPeneliti(w http.ResponseWriter, r *http.Request) {
    var err error

    penelitiID := chi.URLParam(r, "peneliti")

    var peneliti model.Peneliti
    err = routes.App.DB.Preload("Fakultas").Preload("DiciptakanOleh").Where("id = ?", penelitiID).First(&peneliti).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
        } else {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        }
        return
    }

    resp := peneliti.ToResponse()
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) UpdatePeneliti(w http.ResponseWriter, r *http.Request) {
    var err error

    penelitiID := chi.URLParam(r, "peneliti")

    var req resource.PenelitiRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    err = routes.App.DB.Transaction(func (tx *gorm.DB) error {
        var err error

        var peneliti model.Peneliti
        err = routes.App.DB.Where("id = ?", penelitiID).First(&peneliti).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
            } else {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            }
            return api.ErrorHandled
        }

        admin := r.Context().Value("auth_entity").(model.Admin)
        if admin.FakultasID.Valid && (peneliti.FakultasID != admin.FakultasID.String) {
            api.Error{Message: api.ErrForbidden}.Send(w, 403, errors.New("Invalid access rights"))
            return api.ErrorHandled
        }

        peneliti.FromRequest(&req)

        err = routes.App.DB.Save(&peneliti).Error
        if err != nil {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            return api.ErrorHandled
        }

        return nil
    })
    if err != nil {
        if !errors.Is(err, api.ErrorHandled) {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            return
        }
    }

    w.WriteHeader(200)
}

func (routes *RouteCollection) DeletePeneliti(w http.ResponseWriter, r *http.Request) {
    var err error

    penelitiID := chi.URLParam(r, "peneliti")

    err = routes.App.DB.Transaction(func (tx *gorm.DB) error {
        var err error

        var peneliti model.Peneliti
        err = routes.App.DB.Where("id = ?", penelitiID).First(&peneliti).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
            } else {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            }
            return api.ErrorHandled
        }

        admin := r.Context().Value("auth_entity").(model.Admin)
        if admin.FakultasID.Valid && (peneliti.FakultasID != admin.FakultasID.String) {
            api.Error{Message: api.ErrForbidden}.Send(w, 403, errors.New("Invalid access rights"))
            return api.ErrorHandled
        }

        err = routes.App.DB.Delete(&peneliti).Error
        if err != nil {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            return api.ErrorHandled
        }

        return nil
    })
    if err != nil {
        if !errors.Is(err, api.ErrorHandled) {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            return
        }
    }

    w.WriteHeader(200)
}
