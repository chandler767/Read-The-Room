// Created by Chandler Mayo (http://ChandlerMayo.com) and last updated on March 31, 2020.

package main

import (
	"../packages/dontlist"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
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
		// Create a Session with a custom region
		sess := session.Must(session.NewSession(&aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		}))

		// Create a Comprehend client from just a session.
		client := comprehend.New(sess)

		params := comprehend.DetectSentimentInput{}
		params.SetLanguageCode("en")
		params.SetText(string(text))

		req, resp := client.DetectSentimentRequest(&params)

		err := req.Send()
		if err != nil { // resp is now filled
			http.Error(w, "Unable to analyse sentiment.\n", http.StatusBadRequest)
			return
		}

		paramsP := comprehend.DetectKeyPhrasesInput{}
		paramsP.SetLanguageCode("en")
		paramsP.SetText(string(text))

		req2, resp2 := client.DetectKeyPhrasesRequest(&paramsP)

		err = req2.Send()
		if err != nil { // resp2 is now filled
			http.Error(w, "Unable to analyse key phrases.\n", http.StatusBadRequest)
			return
		}

		jsonOffer := map[string]interface{}{ // Forward response to dashboard.
			"Sentiment":  *resp,
			"KeyPhrases": *resp2,
		}
		json.NewEncoder(w).Encode(jsonOffer)
		return
	} else {
		http.Error(w, "Missing text.\n", http.StatusBadRequest)
	}
	return
})

// Opens the specified URL in the default browser of the user.
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
