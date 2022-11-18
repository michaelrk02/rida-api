package service

import (
    "gorm.io/gorm"
)

type CmdParams struct {
    Help bool
    Seed bool
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
