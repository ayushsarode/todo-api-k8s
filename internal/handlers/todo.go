package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/ayushsarode/todo-api-k8s/internal/models"
	"github.com/ayushsarode/todo-api-k8s/internal/store"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	Store *store.TodoStore
}

// NewTodoHandler returns a new TodoHandler with the given TodoStore
func NewTodoHandler(store *store.TodoStore) *TodoHandler {
	return &TodoHandler{Store: store}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, h.Store.List())
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, h.Store.Add(todo))
}

func (h *TodoHandler) CreateTodosBatch(c *gin.Context) {
	var todos []models.Todo
	if err := c.BindJSON(&todos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	results := make([]models.Todo, len(todos))

	for i, t := range todos {
		wg.Add(1)
		go func(i int, t models.Todo) {
			defer wg.Done()
			results[i] = h.Store.Add(t)
		}(i, t)
	}

	wg.Wait()
	c.JSON(http.StatusCreated, results)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updated models.Todo
	if err := c.BindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if todo, ok := h.Store.Update(id, updated); ok {
		c.JSON(http.StatusOK, todo)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
	}
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if h.Store.Delete(id) {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
	}
}
