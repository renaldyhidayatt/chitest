package interfaces

import (
	"chigitaction/models"
	"chigitaction/payload/request"
	"net/http"
)

type ITodoRepository interface {
	Save(todo models.Todo) (models.Todo, error)
	FindAll() ([]models.Todo, error)
	FindByActivityID(ActivityID uint64) ([]models.Todo, error)
	FindOne(id uint64) (models.Todo, error)
	Update(todo models.Todo) (models.Todo, error)
	Delete(todo models.Todo) (bool, error)
}

type ITodoService interface {
	Create(req request.TodoCreateRequest) (models.Todo, error)
	GetAll(ActivityID uint64) ([]models.Todo, error)
	GetOne(id uint64) (models.Todo, error)
	Update(id uint64, req request.TodoUpdateRequest) (models.Todo, error)
	Delete(id uint64) (bool, error)
}

type ITodoHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
