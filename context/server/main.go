package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/josemarjobs/context/log"
)

func main() {
	http.HandleFunc("/", log.Decorate(handler))
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println(ctx, "handler started")
	defer log.Println(ctx, "handler ended")

	fmt.Println("valueFoo: %v", ctx.Value("keyFoo"))

	select {
	case <-time.After(4 * time.Second):
		log.Println(ctx, "handled")
		fmt.Fprintln(w, "Hello World")
	case <-ctx.Done():
		log.Println(ctx, ctx.Err().Error())
		http.Error(w, ctx.Err().Error(), http.StatusInternalServerError)
	}
}
