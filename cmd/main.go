package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"refactored-robot/internal/controller/handlers"
	"refactored-robot/internal/repository"
)

func main() {
	DB := repository.Init()
	h := handlers.New(DB)
	router := mux.NewRouter()

	router.HandleFunc("/users/{id}", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users", h.AddUser).Methods(http.MethodPost)
	router.HandleFunc("/users/get/", h.CheckUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", router)
	log.Println(router)
}
