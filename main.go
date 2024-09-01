package main

import (
	"DemoProofpoint/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.ServePage).Methods("GET")
	r.HandleFunc("/usuarios", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/usuarios", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/usuarios/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/usuarios/{id}", handlers.DeleteUser).Methods("DELETE")

	http.Handle("/", r)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
