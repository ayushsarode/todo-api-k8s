package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayushsarode/todo-api-k8s/internal/handlers"
	"github.com/ayushsarode/todo-api-k8s/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	r := gin.Default()
	todoStore := store.NewTodoStore()
	todoHandler := handlers.NewTodoHandler(todoStore)

	api := r.Group("/api")
	{
		api.GET("/todos", todoHandler.GetTodos)
		api.POST("/todos", todoHandler.CreateTodo)
		api.POST("/todos/batch", todoHandler.CreateTodosBatch)
		api.PUT("/todos/:id", todoHandler.UpdateTodo)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)
	}
	return r
}

func TestCreateTodo(t *testing.T) {
	router := setupTestRouter()

	body := `{"title":"Test Todo"}`

	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBufferString(body))

	req.Header.Set("Content-Type", "applications/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")
}


func TestBatchCreateTodos(t *testing.T) {
	router := setupTestRouter()

	body := `[{"title":"Go"}, {"title":"K8s"}]`
	req, _ := http.NewRequest("POST", "/api/todos/batch", bytes.NewBufferString(body))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Go")
	assert.Contains(t, w.Body.String(), "K8s")
}

func TestGetTodos(t *testing.T) {
	router := setupTestRouter()
	
	req, _ := http.NewRequest("GET", "/api/todos", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
