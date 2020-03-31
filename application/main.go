// Created by Chandler Mayo (http://ChandlerMayo.com) and last updated on March 31, 2020.

package main

import (
	"../packages/dontlist"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

// To run this web application: go get -d && go run main.go

var (
	dashboardTemplate = template.Must(template.ParseFiles("../data/templates/dashboard/dashboard.tmpl"))
	awsAccessKey      string
	awsSecretKey      string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY") // Get our AWS keys from the .env file.
	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
	r := mux.NewRouter()
	r.Handle("/", dashboardHandler).Methods("GET")                                               // Show the dashboard.
	r.Handle("/process/", analysisHandler).Methods("POST")                                       // Process speech with Amazon Comprehend.
	r.PathPrefix("/").Handler(http.FileServer(dontlist.DontListFiles{http.Dir("../data/root")})) // Index file server.
	go open("http://localhost:8091/")
	http.ListenAndServe(":8091", handlers.LoggingHandler(os.Stdout, r))
}

var dashboardHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // http://localhost:8091/  - Show the dashboard.
	w.Header().Set("Content-Type", "text/html")
	err := dashboardTemplate.Execute(w, map[string]interface{}{
		"ComprehendAPI": "http://localhost:8091/process/",
	}) // Show dashboard.
	if err != nil {
		http.Error(w, "Error loading page.", http.StatusInternalServerError)
		return
	}
	return
})

// curl -H "Content-Type: text/plain" -X POST -d 'text to process' http://localhost:8091/process/
var analysisHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // http://localhost:8091/process/
	text, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Missing text to analyse.\n", http.StatusBadRequest)
		return
	}
	if len(text) > 0 {
		w.Write([]byte("Processing text"))
	} else {
		w.Write([]byte("Missing text to analyse.\n"))
	}
	return
})

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
