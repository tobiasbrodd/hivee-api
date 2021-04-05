package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Data struct {
	Content string `json:"content"`
}

var data []Data

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handler(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, data)
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)

	data = append(data, Data{Content: "Some content"})
	router.HandleFunc("/", handler).Methods("GET")

	router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	serverDomain := "localhost"
	serverPort := 8000
	log.Printf("Listening at %s:%d", serverDomain, serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", serverDomain, serverPort), router))
}
