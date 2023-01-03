package router

import (
	"chigitaction/handler"
	"chigitaction/repository"
	"chigitaction/services"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type todosRoute struct {
	handler handler.TodoHandler
	prefix  string
	router  *chi.Mux
	db      *gorm.DB
}

func NewTodoRoutes(prefix string, db *gorm.DB, router *chi.Mux) *todosRoute {
	repository := repository.NewTodoRepository(db)
	service := services.NewTodoService(repository)
	handler := handler.NewTodoHandler(service)

	return &todosRoute{handler: handler, prefix: prefix, router: router, db: db}

}

func (r *todosRoute) TodoRoute() {
	r.router.Route(r.prefix, func(route chi.Router) {
		route.Post("/", r.handler.Create)
		route.Get("/", r.handler.GetAll)
		route.Get("/{id:[0-9]+}", r.handler.GetOne)
		route.Delete("/{id:[0-9]+}", r.handler.Delete)
		route.Put("/{id:[0-9]+}", r.handler.Update)
	})
}
