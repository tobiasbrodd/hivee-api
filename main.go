package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Data struct {
	Content string `json:"content"`
}

var data []Data

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	return r
}

func main() {
	router := router()

	data = append(data, Data{Content: "Some content"})

	http.ListenAndServe(":8000", router)
}
