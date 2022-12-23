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
            api.Error{Message: "Peneliti sudah ada sebelumnya"}.Send(w, 400, err)
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

        err = peneliti.Validate()
        if err != nil {
            api.Error{Message: err.Error()}.Send(w, 400, err)
            return api.ErrorHandled
        }

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
        []string{"peneliti.nidn", "peneliti.nama", "peneliti.h_index", "peneliti.scopus_author_id", "peneliti.gscholar_author_id", "Fakultas.nama", "DiciptakanOleh.nama"},
        []string{"peneliti.nidn", "peneliti.nama", "peneliti.h_index", "Fakultas.nama", "DiciptakanOleh.nama"},
        "peneliti.nama",
    ).FromRequest(r, "")

    owner := func(tx *gorm.DB) *gorm.DB {
        admin := r.Context().Value("auth_entity").(model.Admin)
        if admin.FakultasID.Valid {
            return tx.Where("peneliti.fakultas_id = ?", admin.FakultasID)
        }
        return tx
    }

    var totalItems int64
    err = routes.App.DB.Model(&model.Peneliti{}).Joins("Fakultas").Joins("DiciptakanOleh").Scopes(ds.EnumerationScope).Scopes(owner).Count(&totalItems).Error
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

    penelitiList := []model.Peneliti{}
    err = routes.App.DB.Joins("Fakultas").Joins("DiciptakanOleh").Scopes(ds.PopulationScope).Scopes(owner).Find(&penelitiList).Error
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

        if !admin.FakultasID.Valid {
            if req.FakultasID == nil {
                api.Error{Message: api.ErrMalformedRequest}.Send(w, 400, errors.New("Fakultas must be present"))
                return api.ErrorHandled
            }
            peneliti.FakultasID = *req.FakultasID
        } else {
            peneliti.FakultasID = admin.FakultasID.String
        }

        err = peneliti.Validate()
        if err != nil {
            api.Error{Message: err.Error()}.Send(w, 400, err)
            return api.ErrorHandled
        }

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

func (routes *RouteCollection) GetPenelitiChart(w http.ResponseWriter, r *http.Request) {
    var err error

    fakultasScope := func (tx *gorm.DB) *gorm.DB {
        if r.URL.Query().Get("fakultas_id") != "" {
            return tx.Where("fakultas_id = ?", r.URL.Query().Get("fakultas_id"))
        }
        return tx
    }

    var resp resource.PenelitiChartResponseCollection
    err = routes.App.DB.Model(&model.Peneliti{}).Select("h_index, COUNT(id) AS jumlah").Group("h_index").Scopes(fakultasScope).Find(&resp.Data).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}

func (routes *RouteCollection) GetPenelitiTable(w http.ResponseWriter, r *http.Request) {
    var err error

    var hIndexList []int
    err = routes.App.DB.Model(&model.Peneliti{}).Distinct("h_index").Order("h_index").Pluck("h_index", &hIndexList).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    hIndexMap := make(map[int]int)
    for i, hIndex := range hIndexList {
        hIndexMap[hIndex] = i
    }

    var fakultasList []struct{
        ID string
        Nama string
        Total int
    }
    err = routes.App.DB.Model(&model.Fakultas{}).
        Select("fakultas.id, fakultas.nama, COUNT(peneliti.id) AS total").
        Joins("LEFT JOIN peneliti ON peneliti.fakultas_id = fakultas.id").
        Group("fakultas.id, fakultas.nama").
        Find(&fakultasList).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    fakultasMap := make(map[string]int)
    penelitiCount := 0

    for i, fakultas := range fakultasList {
        fakultasMap[fakultas.ID] = i
        penelitiCount += fakultas.Total
    }

    var resp resource.PenelitiTableResponse

    resp.Headers = make([]string, len(fakultasList))
    resp.Footers = make([]int, len(fakultasList))
    for i, fakultas := range fakultasList {
        resp.Headers[i] = fakultas.Nama
        resp.Footers[i] = fakultas.Total
    }

    resp.Total = penelitiCount

    resp.Rows = make([]resource.PenelitiTableRowResponse, len(hIndexList))
    for i := range resp.Rows {
        row := &resp.Rows[i]

        row.HIndex = hIndexList[i]
        row.Jumlah = 0

        row.Columns = make([]resource.PenelitiTableColumnResponse, len(fakultasList))
        for j := range row.Columns {
            column := &row.Columns[j]

            column.Jumlah = 0
        }
    }

    var result []struct{
        HIndex int
        FakultasID string
        Jumlah int
    }
    err = routes.App.DB.Model(&model.Peneliti{}).
        Select("peneliti.h_index, peneliti.fakultas_id, COUNT(peneliti.id) AS jumlah").
        Joins("JOIN fakultas ON fakultas.id = peneliti.fakultas_id").
        Group("peneliti.h_index, peneliti.fakultas_id").
        Find(&result).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    for _, data := range result {
        row := hIndexMap[data.HIndex]
        column := fakultasMap[data.FakultasID]

        resp.Rows[row].Columns[column].Jumlah = data.Jumlah
        resp.Rows[row].Jumlah += data.Jumlah
    }

    for i := range resp.Rows {
        row := &resp.Rows[i]

        row.Persentase = 0

        for j := range row.Columns {
            column := &row.Columns[j]

            if row.Jumlah != 0 {
                column.Persentase = float64(column.Jumlah) / float64(fakultasList[j].Total) * 100.0
            }

            row.Persentase += column.Persentase
        }
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}
