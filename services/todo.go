package services

import (
	"chigitaction/interfaces"
	"chigitaction/models"
	"chigitaction/payload/request"
	"chigitaction/repository"
	"errors"
	"fmt"
	"time"
)

type TodoService = interfaces.ITodoService

type todoService struct {
	repository repository.TodoRepository
}

func NewTodoService(repository repository.TodoRepository) *todoService {
	return &todoService{repository: repository}
}

func (s *todoService) Create(req request.TodoCreateRequest) (models.Todo, error) {
	todo := models.Todo{
		ActivityGroupID: req.ActivityGroupID,
		Title:           req.Title,
	}

	newTodo, err := s.repository.Save(todo)
	if err != nil {
		return newTodo, err
	}

	return newTodo, err
}

func (s *todoService) GetAll(ActivityID uint64) ([]models.Todo, error) {
	if ActivityID != 0 {
		// Find by activity group id
		todos, err := s.repository.FindByActivityID(ActivityID)
		if err != nil {
			return todos, err
		}
		return todos, nil
	}

	// Find all
	todos, err := s.repository.FindAll()

	if err != nil {
		return todos, err
	}

	return todos, nil
}

func (s *todoService) GetOne(id uint64) (models.Todo, error) {
	// Find all
	todo, err := s.repository.FindOne(id)

	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (s *todoService) Update(id uint64, req request.TodoUpdateRequest) (models.Todo, error) {
	// Find all
	todo, err := s.repository.FindOne(id)
	// If activity group not found
	if todo.ID == 0 {
		message := fmt.Sprintf("Todo with ID %d Not Found", id)
		return todo, errors.New(message)
	}

	if err != nil {
		return todo, err
	}

	// Change field title
	if req.Title != "" {
		todo.Title = req.Title
	}
	// Change field is active, if value req.IsActive is true
	if !req.IsActive {
		todo.IsActive = req.IsActive
	} else {
		todo.IsActive = true

	}

	todo.UpdatedAt = time.Now()

	// Update
	updatedTodo, err := s.repository.Update(todo)
	if err != nil {
		return updatedTodo, err
	}

	return updatedTodo, nil
}

func (s *todoService) Delete(id uint64) (bool, error) {
	// Find one
	todo, err := s.repository.FindOne(id)
	// If activity group not found
	if todo.ID == 0 {
		message := fmt.Sprintf("Todo with ID %d Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	ok, err := s.repository.Delete(todo)
	if err != nil {
		return false, err
	}

	return ok, nil
}
