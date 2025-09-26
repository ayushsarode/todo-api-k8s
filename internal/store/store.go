package store

import (
	"sync"

	"github.com/ayushsarode/todo-api-k8s/internal/models"
)

type TodoStore struct {
	sync.RWMutex
	todos  []models.Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{todos: []models.Todo{}, nextID: 1}
}

func (s *TodoStore) Add(todo models.Todo) models.Todo {
	s.Lock()
	defer s.Unlock()
	
	todo.ID = s.nextID
	s.nextID++
	s.todos = append(s.todos, todo)
	return todo
}

func (s *TodoStore) List() []models.Todo {
	s.RLock()
	defer s.RUnlock()
	
	// Return a copy to avoid concurrent access issues
	result := make([]models.Todo, len(s.todos))
	copy(result, s.todos)
	return result
}

func (s *TodoStore) Update(id int, updated models.Todo) (models.Todo, bool) {
	s.Lock()
	defer s.Unlock()
	
	for i, t := range s.todos {
		if t.ID == id {
			s.todos[i].Title = updated.Title
			s.todos[i].Done = updated.Done
			return s.todos[i], true
		}
	}
	return models.Todo{}, false
}

func (s *TodoStore) Delete(id int) bool {
	s.Lock()
	defer s.Unlock()
	
	for i, t := range s.todos {
		if t.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return true
		}
	}
	return false
}
