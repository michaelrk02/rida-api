package route

import (
    "github.com/go-chi/chi/v5"
    "github.com/michaelrk02/rida-api/middleware"
    "github.com/michaelrk02/rida-api/service"
)

func InitRoutesV1(r chi.Router, app *service.Application) {
    r.Post("/test", app.RouteV1.CreateTest)
    r.Get("/test", app.RouteV1.GetAllTest)
    r.Get("/test/{test}", app.RouteV1.GetTest)
    r.Put("/test/{test}", app.RouteV1.UpdateTest)
    r.Delete("/test/{test}", app.RouteV1.DeleteTest)

    r.Post("/admin/login", app.RouteV1.LoginAdmin)
    r.Group(func (r chi.Router) {
        r.Use(middleware.AuthAdmin(true))
        r.Post("/admin", app.RouteV1.CreateAdmin)
        r.Get("/admin", app.RouteV1.GetAllAdmin)
        r.Delete("/admin/{admin}", app.RouteV1.DeleteAdmin)
    })
    r.Group(func (r chi.Router) {
        r.Use(middleware.AuthAdmin(false))
        r.Get("/admin/{admin}", app.RouteV1.GetAdmin)
        r.Put("/admin/{admin}", app.RouteV1.UpdateAdmin)
    })

    r.Group(func (r chi.Router) {
        r.Use(middleware.AuthAdmin(false))
        r.Post("/peneliti", app.RouteV1.CreatePeneliti)
        r.Get("/peneliti", app.RouteV1.GetAllPeneliti)
        r.Get("/peneliti/{peneliti}", app.RouteV1.GetPeneliti)
        r.Put("/peneliti/{peneliti}", app.RouteV1.UpdatePeneliti)
        r.Delete("/peneliti", app.RouteV1.DeletePeneliti)
    })
}
