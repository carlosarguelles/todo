package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/carlosarguelles/todo/internal/infra/cli"
	"github.com/carlosarguelles/todo/internal/infra/db"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  add <note> - Add a new note")
		fmt.Println("  list       - List all notes")
		fmt.Println("  delete <id> - Delete a note by ID")
		return
	}

	key := os.Getenv("TODO_KEY")

	if key == "" {
		fmt.Println("error: TODO_KEY not set")
		os.Exit(1)
	}

	repo := db.NewRedisNodeRepository(rdb, key)
	cliTodo := cli.NewTodoCli(repo)

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a note to add.")
			return
		}
		note := strings.Join(os.Args[2:], " ")
		cliTodo.AddNote(ctx, note)

	case "list":
		cliTodo.ListNotes(ctx)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide the ID of the note to delete.")
			return
		}
		id := os.Args[2]
		cliTodo.DeleteNote(ctx, id)

	default:
		fmt.Println("Unknown command:", command)
	}
}
