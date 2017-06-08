package log

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type key int

const requestIDKey = key(42)

func Println(ctx context.Context, msg string) {
	id, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		log.Println("could not find request ID in context")
		return
	}
	log.Printf("[%d] - %s\n", id, msg)
}

func Decorate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().Unix())
		ctx := r.Context()
		id := rand.Int63()
		ctx = context.WithValue(ctx, requestIDKey, id)
		f(w, r.WithContext(ctx))
	}
}
