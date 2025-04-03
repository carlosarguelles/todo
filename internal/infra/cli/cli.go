package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/carlosarguelles/todo/internal/app"
)

type TodoCli struct {
	repository app.NoteRepository
}

func NewTodoCli(repository app.NoteRepository) *TodoCli {
	return &TodoCli{repository}
}

func (c *TodoCli) AddNote(ctx context.Context, note string) {
	err := c.repository.AddNote(ctx, note)
	if err != nil {
		fmt.Println("failed to add note")
		os.Exit(1)
	}
	fmt.Println("note added")
}

func (c *TodoCli) ListNotes(ctx context.Context) {
	notes, err := c.repository.GetAllNotes(ctx)
	if err != nil {
		fmt.Println("failed to fetch notes")
		os.Exit(1)
	}

	if len(notes) == 0 {
		fmt.Println("no notes found")
		return
	}

	for _, note := range notes {
		fmt.Printf("%s: %s\n", note.ID, note.Text)
	}
}

func (c *TodoCli) DeleteNote(ctx context.Context, id string) {
	err := c.repository.DeleteNote(ctx, id)
	if err != nil {
		fmt.Println("failed to delete note")
		os.Exit(1)
	}
	fmt.Printf("note with ID %s deleted.\n", id)
}
