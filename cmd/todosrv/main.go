package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/carlosarguelles/todo/internal/infra/api"
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
	port := flag.String("port", ":8080", "the port number for the server")
	key := flag.String("key", "todo", "the key for the notes")
	flag.Parse()

	repo := db.NewRedisNodeRepository(rdb, *key)
	todoApi := api.NewTodoApi(repo)

	http.HandleFunc("/add", todoApi.Add)
	http.HandleFunc("/list", todoApi.List)
	http.HandleFunc("/delete", todoApi.Delete)

	fmt.Printf("Starting server on port %s\n", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
