package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasbrodd/hivee-api/internal/storage"
	"gopkg.in/yaml.v3"
)

type storageServer struct {
	store *storage.Storage
}

func NewStorageServer(c config) *storageServer {
	store := storage.New(c.Influx.Token, c.Influx.Host, c.Influx.Port, "Hivee")
	return &storageServer{store: store}
}

type config struct {
	Influx struct {
		Token string `yaml:"token"`
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
	}
}

func (c *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Errorf("Config: %v", err.Error())
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Errorf("Config: %v", err.Error())
	}

	return c
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *storageServer) getMeasureHistoryHandler(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	measurement := v.Get("measurement")
	location := v.Get("location")
	history := s.store.ReadMeasureHistory(measurement, location)

	renderJSON(w, history)
}

func main() {
	formatter := &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true}
	log.SetFormatter(formatter)

	var c config
	c.getConfig()

	server := NewStorageServer(c)
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/history", server.getMeasureHistoryHandler).Methods("GET")

	router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	serverDomain := "0.0.0.0"
	serverPort := 8000
	log.Infof("Listening at %s:%d", serverDomain, serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", serverDomain, serverPort), router))
	server.store.Close()
}
