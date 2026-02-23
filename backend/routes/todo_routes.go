package routes

import (
	"todo-fullstack/backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTodoRoutes(r *gin.Engine, todoController controllers.TodoController) {
	api := r.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.GET("/", todoController.GetAllTodos)
			todos.GET("/:id", todoController.GetTodoByID)
			todos.POST("/", todoController.CreateTodo)
			todos.PUT("/:id", todoController.UpdateTodo)
			todos.DELETE("/:id", todoController.DeleteTodo)
		}
	}
}
