package v1

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/golang-jwt/jwt/v4"
    "github.com/michaelrk02/rida-api/api"
    "github.com/michaelrk02/rida-api/model"
    "github.com/michaelrk02/rida-api/resource"
    "github.com/michaelrk02/rida-api/service"
    "gorm.io/gorm"
)

func (routes *RouteCollection) LoginAdmin(w http.ResponseWriter, r *http.Request) {
    var err error

    var req resource.AdminRequestOnLogin
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    var admin model.Admin
    err = routes.App.DB.Where("email = ?", req.Email).First(&admin).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            api.Error{Message: api.ErrUnauthorized}.Send(w, 401, err)
        } else {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        }
        return
    }

    passwordHash := sha256.New()
    passwordHash.Write([]byte(req.Password))
    if hex.EncodeToString(passwordHash.Sum(nil)) != admin.Password {
        api.Error{Message: api.ErrUnauthorized}.Send(w, 401, errors.New("Invalid password"))
        return
    }

    tokenDuration, _ := time.ParseDuration("24h")
    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": admin.ID,
        "exp": time.Now().Add(tokenDuration).Unix(),
    }).SignedString([]byte(os.Getenv("APP_KEY")))
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    role := ""
    if admin.FakultasID.Valid {
        role = "admin"
    } else {
        role = "superadmin"
    }

    resp := resource.AdminResponseOnLogin{
        ID: admin.ID,
        Role: role,
        Token: token,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) CreateAdmin(w http.ResponseWriter, r *http.Request) {
    var err error

    var req resource.AdminRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    var resp resource.AdminResponseOnCreate

    err = routes.App.DB.Transaction(func(tx *gorm.DB) error {
        var err error

        var admin model.Admin

        err = tx.Where("email = ?", req.Email).First(&admin).Error
        if err == nil {
            api.Error{Message: "Admin already exists"}.Send(w, 400, err)
            return api.ErrorHandled
        }
        if err != nil {
            if !errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
                return api.ErrorHandled
            }
        }

        admin = model.Admin{}
        admin.FromRequest(&req)

        err = admin.Validate(model.ADMIN_VALIDATE_ALL)
        if err != nil {
            api.Error{Message: err.Error()}.Send(w, 400, err)
            return api.ErrorHandled
        }

        admin.HashPassword()

        err = routes.App.DB.Create(&admin).Error
        if err != nil {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            return api.ErrorHandled
        }

        resp = admin.ToResponseOnCreate()

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

func (routes *RouteCollection) GetAllAdmin(w http.ResponseWriter, r *http.Request) {
    var err error

    ds := service.NewDataSource(
        100,
        []string{"admin.nama", "admin.email", "Fakultas.nama"},
        []string{"admin.nama", "admin.email", "Fakultas.nama"},
        "admin.nama",
    ).FromRequest(r, "")

    roleScope := func (tx *gorm.DB) *gorm.DB {
        return tx.Where("fakultas_id IS NOT NULL")
    }

    var totalItems int64
    err = routes.App.DB.Model(&model.Admin{}).Joins("Fakultas").Scopes(ds.EnumerationScope).Scopes(roleScope).Count(&totalItems).Error
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

    var adminList []model.Admin
    err = routes.App.DB.Joins("Fakultas").Scopes(ds.PopulationScope).Scopes(roleScope).Find(&adminList).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    var resp resource.AdminResponseCollection
    resp.Population = ds.TotalItems
    resp.Display = ds.Display
    resp.Page = ds.Page
    resp.MaxPage = ds.MaxPage
    resp.Data = make([]resource.AdminResponse, len(adminList))
    for i := range adminList {
        resp.Data[i] = adminList[i].ToResponse()
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) GetAdmin(w http.ResponseWriter, r *http.Request) {
    var err error

    adminID := chi.URLParam(r, "admin")

    requestAdminID := r.Context().Value("auth_id").(string)
    isSuperadmin := r.Context().Value("auth_is_superadmin").(bool)
    if adminID != requestAdminID && !isSuperadmin {
        api.Error{Message: api.ErrUnauthorized}.Send(w, 401, err)
        return
    }

    var admin model.Admin
    err = routes.App.DB.Preload("Fakultas").Where("id = ?", adminID).First(&admin).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
        } else {
            api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        }
        return
    }

    resp := admin.ToResponse()
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
    var err error

    adminID := chi.URLParam(r, "admin")

    requestAdminID := r.Context().Value("auth_id").(string)
    isSuperadmin := r.Context().Value("auth_is_superadmin").(bool)
    if adminID != requestAdminID && !isSuperadmin {
        api.Error{Message: api.ErrUnauthorized}.Send(w, 401, err)
        return
    }

    var req resource.AdminRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    err = routes.App.DB.Transaction(func (tx *gorm.DB) error {
        var err error

        var admin model.Admin
        err = routes.App.DB.Where("id = ?", adminID).First(&admin).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
            } else {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            }
            return api.ErrorHandled
        }

        admin.FromRequest(&req).AssignID(adminID)

        err = admin.Validate(model.ADMIN_VALIDATE_ALL)
        if err != nil {
            api.Error{Message: err.Error()}.Send(w, 400, err)
            return api.ErrorHandled
        }

        admin.HashPassword()

        err = routes.App.DB.Save(&admin).Error
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

func (routes *RouteCollection) UpdateAdminPassword(w http.ResponseWriter, r *http.Request) {
    var err error

    adminID := chi.URLParam(r, "admin")
    requestAdminID := r.Context().Value("auth_id").(string)
    if adminID != requestAdminID {
        api.Error{Message: api.ErrForbidden}.Send(w, 403, errors.New("Invalid access"))
        return
    }

    var req resource.AdminRequestOnUpdatePassword
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, err)
        return
    }

    err = routes.App.DB.Transaction(func (tx *gorm.DB) error {
        var err error

        var admin model.Admin
        err = routes.App.DB.Where("id = ?", adminID).First(&admin).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
            } else {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            }
            return api.ErrorHandled
        }

        passwordHash := sha256.New()
        passwordHash.Write([]byte(req.OldPassword))
        if hex.EncodeToString(passwordHash.Sum(nil)) != admin.Password {
            api.Error{Message: "Password yang lama tidak valid"}.Send(w, 401, errors.New(api.ErrUnauthorized))
            return api.ErrorHandled
        }

        admin.Password = req.NewPassword
        err = admin.Validate(model.ADMIN_VALIDATE_PASSWORD)
        if err != nil {
            api.Error{Message: err.Error()}.Send(w, 400, err)
            return api.ErrorHandled
        }

        admin.HashPassword()

        err = routes.App.DB.Save(&admin).Error
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

func (routes *RouteCollection) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
    var err error

    adminID := chi.URLParam(r, "admin")

    err = routes.App.DB.Transaction(func (tx *gorm.DB) error {
        var err error

        var admin model.Admin
        err = routes.App.DB.Where("id = ?", adminID).First(&admin).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                api.Error{Message: api.ErrNotFound}.Send(w, 404, err)
            } else {
                api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
            }
            return api.ErrorHandled
        }

        err = routes.App.DB.Delete(&admin).Error
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
