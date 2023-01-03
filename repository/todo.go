package repository

import (
	"chigitaction/interfaces"
	"chigitaction/models"

	"gorm.io/gorm"
)

type TodoRepository = interfaces.ITodoRepository

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *todoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll() ([]models.Todo, error) {
	var todos []models.Todo

	err := r.db.Find(&todos).Error
	if err != nil {
		return todos, nil
	}

	return todos, nil
}

func (r *todoRepository) FindByActivityID(ActivityID uint64) ([]models.Todo, error) {
	var todos []models.Todo

	err := r.db.Where("activity_group_id = ?", ActivityID).Find(&todos).Error
	if err != nil {
		return todos, nil
	}

	return todos, nil
}

func (r *todoRepository) FindOne(id uint64) (models.Todo, error) {
	var todo models.Todo

	err := r.db.Where("id = ?", id).Find(&todo).Error
	if err != nil {
		return todo, nil
	}

	return todo, nil
}

func (r *todoRepository) Save(todo models.Todo) (models.Todo, error) {
	err := r.db.Create(&todo).Error
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (r *todoRepository) Update(todo models.Todo) (models.Todo, error) {
	err := r.db.Save(&todo).Error
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (r *todoRepository) Delete(todo models.Todo) (bool, error) {
	err := r.db.Delete(&todo).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
