package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NanobyteRuata/go-taskmanager/internal/models"
	"github.com/gorilla/mux"
)

type Handler struct {
	store models.TaskStore
}

func NewHandler(store models.TaskStore) *Handler {
	return &Handler{
		store: store,
	}
}

// Router returns a configured router for the API
func (h *Handler) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", h.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", h.CreateTask).Methods("POST")

	return r
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		DueDate string `json:"due_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	task := models.NewTask(req.Title)

	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid due date format. Use YYYY-MM-DD")
			return
		}

		task.DueDate = dueDate

		savedTask, err := h.store.Create(task)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create task")
			return
		}

		respondWithJSON(w, http.StatusCreated, savedTask)
	}
}
