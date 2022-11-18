package v1

import (
    "github.com/michaelrk02/rida-api/service"
)

type RouteCollection struct {
    service.RouteCollectionV1

    App *service.Application
}
