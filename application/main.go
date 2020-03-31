// Created by Chandler Mayo (http://ChandlerMayo.com) and last updated on March 31, 2020.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

// To run this web application: go run main.go

var (
	dashboardTemplate = template.Must(template.New("dashboard", Asset).Parse("../templates/dashboard.tmpl"))
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.Handle("/", dashboardHandler).Methods("GET")         // Show the dashboard.
	r.Handle("/process/", analysisHandler).Methods("POST") // Process speech with Amazon Comprehend.
	http.ListenAndServe(":8091", handlers.LoggingHandler(os.Stdout, r))
}

var dashboardHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // http://localhost:8091/  - Check the status of the API.
	w.Write([]byte("This is your API. It is online.\n"))
})

// curl -H "Content-Type: text/plain" -X POST -d 'text to analyse' http://localhost:8091/process/
var analysisHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // http://localhost:8091/process/
	text, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Missing text to analyse.\n", http.StatusBadRequest)
		return
	}
	if len(text) > 0 {
		w.Write([]byte("Your name is: " + name + ". This is a post request.\n"))
	} else {
		w.Write([]byte("Missing text to analyse.\n"))
	}
})
