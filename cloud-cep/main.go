package main

import (
	"log"
	"net/http"

	"jhonathannc/pos-go-desafios/cloud-cep/handler"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", handler.RootHandler).Methods("GET")
	r.HandleFunc("/{cep}", handler.WeatherByCEPHandler).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
