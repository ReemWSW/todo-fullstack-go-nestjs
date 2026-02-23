package repository

import (
	"todo-fullstack/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepository interface {
	FindAll() ([]models.Todo, error)
	FindByID(id uuid.UUID) (models.Todo, error)
	Create(todo models.Todo) (models.Todo, error)
	Update(todo models.Todo) (models.Todo, error)
	Delete(id uuid.UUID) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll() ([]models.Todo, error) {
	var todos []models.Todo
	if err := r.db.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) FindByID(id uuid.UUID) (models.Todo, error) {
	var todo models.Todo
	if err := r.db.First(&todo, "id = ?", id).Error; err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *todoRepository) Create(todo models.Todo) (models.Todo, error) {
	if err := r.db.Create(&todo).Error; err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *todoRepository) Update(todo models.Todo) (models.Todo, error) {
	if err := r.db.Save(&todo).Error; err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *todoRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.Todo{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
