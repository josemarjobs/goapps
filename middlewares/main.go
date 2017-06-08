package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bmizerany/pat"
	"github.com/urfave/negroni"
)

func main() {
	doneCh := waitForTemination()
	go func() {
		<-doneCh
		log.Println("exiting...")
		os.Exit(0)
	}()

	m := pat.New()
	m.Get("/", http.HandlerFunc(HomeHandler))
	m.Get("/me", http.HandlerFunc(ProfileHandler))

	n := negroni.Classic()
	n.UseHandler(m)
	n.Use(negroni.HandlerFunc(MyMiddleware))
	http.ListenAndServe(":3000", n)
}

func MyMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Before call")
	next(w, r)
	log.Println("After call")
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Authenticating...")
	next(w, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Profile stuff")
}
func waitForTemination() <-chan struct{} {
	doneCh := make(chan struct{}, 1)
	ch := make(chan os.Signal, 1)
	go func() {
		<-ch
		log.Println("Stopping...")
		doneCh <- struct{}{}
	}()
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	return doneCh
}
