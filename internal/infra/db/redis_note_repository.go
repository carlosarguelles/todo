package db

import (
	"context"
	"fmt"

	"github.com/carlosarguelles/todo/internal/dom"
	"github.com/go-redis/redis/v8"
)

type RedisNodeRepository struct {
	cmd redis.Cmdable
	key string
}

func NewRedisNodeRepository(cmd redis.Cmdable, key string) *RedisNodeRepository {
	return &RedisNodeRepository{cmd, key}
}

func (r *RedisNodeRepository) AddNote(ctx context.Context, note string) error {
	id, err := r.cmd.Incr(ctx, fmt.Sprintf("%s:id", r.key)).Result()
	if err != nil {
		return err
	}
	err = r.cmd.Set(ctx, fmt.Sprintf("%s:%d", r.key, id), note, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisNodeRepository) GetAllNotes(ctx context.Context) ([]dom.Note, error) {
	keys, err := r.cmd.Keys(ctx, fmt.Sprintf("%s:*", r.key)).Result()
	if err != nil {
		return nil, err
	}
	var notes []dom.Note
	for _, key := range keys {
		if key == "id" {
			continue
		}
		note, err := r.cmd.Get(ctx, key).Result()
		if err != nil {
			fmt.Printf("failed to get note %s: %v\n", key, err)
			continue
		}
		notes = append(notes, dom.Note{
			ID:   key[len(r.key+":"):],
			Text: note,
		})
	}
	return notes, nil
}

func (r *RedisNodeRepository) DeleteNote(ctx context.Context, id string) error {
	if id == "id" {
		return nil
	}
	key := fmt.Sprintf("%s:%s", r.key, id)
	err := r.cmd.Unlink(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
