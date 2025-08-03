package routes

import "github.com/go-redis/redis/v8"

// Handler holds shared dependencies
type Handler struct {
	Rdb0 *redis.Client
	Rdb1 *redis.Client
}

// NewHandler wires up Redis clients into a single struct.
func NewHandler(rdb0, rdb1 *redis.Client) *Handler {
	return &Handler{
		Rdb0: rdb0,
		Rdb1: rdb1,
	}
}
