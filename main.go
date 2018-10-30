package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"encoding/json"
)

// our main function
func main() {
	router := mux.NewRouter()
	log.Print("Listening on port 8080")
	router.HandleFunc("/people", GetPeople).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}


func GetPeople(w http.ResponseWriter, r *http.Request) {
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe"}	)

	json.NewEncoder(w).Encode(people)
}

type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
}

var people []Person