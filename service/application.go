package service

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"

    "github.com/michaelrk02/rida-api/model"
    "gorm.io/gorm"
)

type CmdParams struct {
    Help bool
    SeedLocal bool
    SeedRemote bool
    Sync bool
    Daemon bool
    Panic bool
}

type Application struct {
    Params CmdParams
    DB *gorm.DB

    RouteV1 RouteCollectionV1
}

var appInstance *Application

func InitApplication() *Application {
    appInstance = &Application{}

    return appInstance
}

func GetApplication() *Application {
    return appInstance
}

func (app *Application) SyncHIndex() {
    var err error

    penelitiList := []model.Peneliti{}
    err = app.DB.Where("is_remote").Find(&penelitiList).Error
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
        return
    }

    for i, peneliti := range penelitiList {
        fmt.Printf(" [%d / %d] %s ... ", i + 1, len(penelitiList), peneliti.Nama)

        resp, err := http.Get(fmt.Sprintf(
            "https://serpapi.com/search.json?%s",
            url.Values{
                "engine": []string{"google_scholar_author"},
                "author_id": []string{peneliti.GscholarAuthorID},
                "api_key": []string{os.Getenv("SERPAPI_KEY")},
            }.Encode(),
        ))
        if err != nil {
            fmt.Println("error")
            continue
        }

        var data struct {
            Error string `json:"error"`
            CitedBy struct {
                Table []map[string]struct{
                    All int `json:"all"`
                } `json:"table"`
            } `json:"cited_by"`
        }
        err = json.NewDecoder(resp.Body).Decode(&data)
        resp.Body.Close()
        if err != nil {
            fmt.Println("error")
            continue
        }

        if data.Error != "" {
            fmt.Println("error")
            continue
        }

        for _, table := range data.CitedBy.Table {
            if hIndex, ok := table["h_index"]; ok {
                peneliti.HIndex = hIndex.All
                fmt.Printf("(%d) ", peneliti.HIndex)
                break
            }
        }

        err = app.DB.Save(&peneliti).Error
        if err != nil {
            fmt.Println("error")
            continue
        }

        fmt.Println("ok")
    }
}
