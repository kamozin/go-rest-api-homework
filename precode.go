package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Handler for getting all tasks
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	// Marshal tasks to JSON
	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJSON)
}

// Handler for adding a task
func addTask(w http.ResponseWriter, r *http.Request) {
	// Decode JSON request body to Task struct
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Add task to map
	tasks[newTask.ID] = newTask
	// Return success response
	w.WriteHeader(http.StatusCreated)
}

// Handler for getting a task by ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	// Get task ID from URL parameters
	taskID := chi.URLParam(r, "id")
	// Check if task exists
	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}
	// Marshal task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")
	// Write JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(taskJSON)
}

// Handler for deleting a task by ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	// Get task ID from URL parameters
	taskID := chi.URLParam(r, "id")
	// Check if task exists
	_, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}
	// Delete task from map
	delete(tasks, taskID)
	// Return success response
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Registering handlers
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
