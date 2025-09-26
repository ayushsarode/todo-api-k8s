package main

import (
	"github.com/ayushsarode/todo-api-k8s/internal/handlers"
	"github.com/ayushsarode/todo-api-k8s/internal/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	todoHandler := &handlers.TodoHandler{Store: store.NewTodoStore()}

	api := r.Group("/api")
	{
		api.GET("/todos", todoHandler.GetTodos)
		api.POST("/todos", todoHandler.CreateTodo)
		api.POST("/todos/batch", todoHandler.CreateTodosBatch)
		api.PUT("/todos/:id", todoHandler.UpdateTodo)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)
	}

	r.GET("/healthz", func(c *gin.Context)  {
		c.String(http.StatusOK, "ok")
	})

	r.Run(":8080")
}

