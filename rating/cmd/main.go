package main

import (
	"log"
	"net/http"

	"movieapp.com/rating/internal/controller/rating"
	httphandler "movieapp.com/rating/internal/handler/http"
	"movieapp.com/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating server")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}