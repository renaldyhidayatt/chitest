package router

import (
	"chigitaction/handler"
	"chigitaction/repository"
	"chigitaction/services"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type activityRoute struct {
	handler handler.ActivityHandler
	prefix  string
	router  *chi.Mux
	db      *gorm.DB
}

func NewActivityRoutes(prefix string, db *gorm.DB, router *chi.Mux) *activityRoute {
	repository := repository.NewRepositoryActivity(db)
	service := services.NewActivityService(repository)
	handler := handler.NewActivityHandler(service)

	return &activityRoute{handler: handler, prefix: prefix, router: router, db: db}
}

func (r *activityRoute) ActivityRoute() {
	r.router.Route(r.prefix, func(route chi.Router) {
		route.Post("/", r.handler.Create)
		route.Get("/", r.handler.GetAll)
		route.Get("/{id:[0-9]+}", r.handler.GetOne)
		route.Delete("/{id:[0-9]+}", r.handler.Delete)
		route.Put("/{id:[0-9]+}", r.handler.Update)
	})
}
