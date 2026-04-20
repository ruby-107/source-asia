package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ruby-107/source-asia/internal/handler"
	"github.com/ruby-107/source-asia/internal/repository"
	"github.com/ruby-107/source-asia/internal/service"
)

func main() {
	godotenv.Load()

	db, err := repository.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	go service.StartWorker(db)

	r := mux.NewRouter()

	limiter := service.NewRateLimiter()
	h := handler.NewHandler(limiter, db)

	r.HandleFunc("/request", h.Request).Methods("POST")
	r.HandleFunc("/stats", h.Stats).Methods("GET")

	log.Println("Server running on :" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
