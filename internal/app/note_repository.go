package app

import (
	"context"

	"github.com/carlosarguelles/todo/internal/dom"
)

type NoteRepository interface {
	AddNote(context.Context, string) error
	GetAllNotes(context.Context) ([]dom.Note, error)
	DeleteNote(context.Context, string) error
}
