package v1

import (
    "encoding/json"
    "net/http"

    "github.com/michaelrk02/rida-api/api"
    "github.com/michaelrk02/rida-api/model"
    "github.com/michaelrk02/rida-api/resource"
)

func (routes *RouteCollection) GetAllFakultas(w http.ResponseWriter, r *http.Request) {
    var err error

    var fakultasList []model.Fakultas
    err = routes.App.DB.Find(&fakultasList).Error
    if err != nil {
        api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
        return
    }

    var resp resource.FakultasResponseCollection
    resp.Count = len(fakultasList)
    resp.Data = make([]resource.FakultasResponse, len(fakultasList))
    for i := range fakultasList {
        resp.Data[i] = fakultasList[i].ToResponse()
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(&resp)
}
