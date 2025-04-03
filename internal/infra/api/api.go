package api

import (
	"encoding/json"
	"net/http"

	"github.com/carlosarguelles/todo/internal/app"
)

type TodoApi struct {
	repository app.NoteRepository
}

func NewTodoApi(repository app.NoteRepository) *TodoApi {
	return &TodoApi{repository}
}

func (a *TodoApi) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := r.URL.Query().Get("text")
	if note == "" {
		http.Error(w, "Note text is required", http.StatusBadRequest)
		return
	}

	err := a.repository.AddNote(r.Context(), note)
	if err != nil {
		http.Error(w, "Failed to add note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *TodoApi) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	notes, err := a.repository.GetAllNotes(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (a *TodoApi) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}
	err := a.repository.DeleteNote(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
