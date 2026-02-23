package controllers

import (
	"net/http"

	"todo-fullstack/backend/models"
	"todo-fullstack/backend/services"
	"todo-fullstack/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TodoController interface {
	GetAllTodos(c *gin.Context)
	GetTodoByID(c *gin.Context)
	CreateTodo(c *gin.Context)
	UpdateTodo(c *gin.Context)
	DeleteTodo(c *gin.Context)
}

type todoController struct {
	todoService services.TodoService
	validate    *validator.Validate
}

func NewTodoController(service services.TodoService) TodoController {
	return &todoController{
		todoService: service,
		validate:    validator.New(),
	}
}

func (ctrl *todoController) GetAllTodos(c *gin.Context) {
	todos, err := ctrl.todoService.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(todos, "Todos retrieved successfully"))
}

func (ctrl *todoController) GetTodoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Todo ID"))
		return
	}

	todo, err := ctrl.todoService.GetTodoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Todo not found"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(todo, "Todo retrieved successfully"))
}

func (ctrl *todoController) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	// Validate the struct
	if err := ctrl.validate.Struct(todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	createdTodo, err := ctrl.todoService.CreateTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse(createdTodo, "Todo created successfully"))
}

func (ctrl *todoController) UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Todo ID"))
		return
	}

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	// Validate the struct
	if err := ctrl.validate.Struct(todo); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	updatedTodo, err := ctrl.todoService.UpdateTodo(id, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(updatedTodo, "Todo updated successfully"))
}

func (ctrl *todoController) DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Todo ID"))
		return
	}

	if err := ctrl.todoService.DeleteTodo(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Todo deleted successfully"))
}
