package services

import (
	"errors"
	"todo-fullstack/backend/models"
	"todo-fullstack/backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoService interface {
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id uuid.UUID) (models.Todo, error)
	CreateTodo(todo models.Todo) (models.Todo, error)
	UpdateTodo(id uuid.UUID, todo models.Todo) (models.Todo, error)
	DeleteTodo(id uuid.UUID) error
}

type todoService struct {
	todoRepository repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{todoRepository: repo}
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.todoRepository.FindAll()
}

func (s *todoService) GetTodoByID(id uuid.UUID) (models.Todo, error) {
	return s.todoRepository.FindByID(id)
}

func (s *todoService) CreateTodo(todo models.Todo) (models.Todo, error) {
	return s.todoRepository.Create(todo)
}

func (s *todoService) UpdateTodo(id uuid.UUID, todo models.Todo) (models.Todo, error) {
	existingTodo, err := s.todoRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Todo{}, errors.New("Todo not found")
		}
		return models.Todo{}, err
	}

	existingTodo.Title = todo.Title
	existingTodo.Description = todo.Description
	existingTodo.Completed = todo.Completed

	return s.todoRepository.Update(existingTodo)
}

func (s *todoService) DeleteTodo(id uuid.UUID) error {
	_, err := s.todoRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Todo not found")
		}
		return err
	}
	return s.todoRepository.Delete(id)
}
