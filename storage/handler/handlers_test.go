package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"telegram_reminder_bot/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockService struct {
	CreateTaskFn func(task models.Task) error
	TasksFn      func() ([]models.Task, error)
}

func (m *MockService) CreateTask(task models.Task) error {
	// Например, просто возвращаем nil, т.е. успешно
	return nil
}

func (m *MockService) Tasks() ([]models.Task, error) {
	// Возвращаем пустой список задач
	return []models.Task{}, nil
}

func TestCreateTask(t *testing.T) {
	mockService := &MockService{}

	req, err := http.NewRequest("POST", "/tasks", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := mockService.CreateTask(models.Task{})
		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetTasks(t *testing.T) {
	mockService := &MockService{}

	req, err := http.NewRequest("GET", "/tasks", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tasks, err := mockService.Tasks()
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	}).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}
