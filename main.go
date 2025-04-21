package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"regexp"
	"strings"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}


type ReleaseResponse struct {
	JetId							string `json:"jetId"`
	ReleaseReady					bool `json:"releaseReady"`
	JetConsoleUrl					string `json:"jetConsoleUrl"`
	ReleaseReadyMessage				[]string `json:"releaseReadyMessage"`
}


type ServiceNowResponse struct {
	State 							string `json:"state"`
	StartTime 						string `json:"startTime"`
	EndTime							string `json:"endTime"`
	MainConfigurationItem			MainConfigurationItem `json:"mainConfigurationItem"`
}

type JFrog struct {
	Created				string `json:"created"`
	Checksums			map[string]string `json:"checksums"`
	DownloadUri			string `json:"downloadUri"`
	Garbage 			string `json:"garbage"`
}

type ItemProperties struct{

	JetId				string `json:"jet-id"`
	SealId				string `json:"seal_Id"`
	ProjectName 		string `json:"project_key"`
}


type MainConfigurationItem struct {
	Name 							string `json:"name"`
	Number 							ConfigurationItemNumber `json:"number"`
}

type ConfigurationItemNumber struct {
	Identifier						string `json:"identifier"`
}

// jsonHandler writes a JSON response to the client.
func evidenceCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of Response.
	releaseResponse := ReleaseResponse{
		JetId : "sample jet id",
		ReleaseReady: true,
		JetConsoleUrl : "[console uri link]",
		ReleaseReadyMessage: []string{"This is a mock message"},
	}

	// Marshal the struct to JSON.
	jsonData, err := json.Marshal(releaseResponse)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate header and status code.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func snowCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of Response.
	releaseResponse := ServiceNowResponse{
		State 							: "Scheduled",
		StartTime 						: "2025-02-23T05:25:10Z",
		EndTime							: "2025-03-05T05:25:10Z",
		MainConfigurationItem			: MainConfigurationItem{
			Name : "gerf",
			Number : ConfigurationItemNumber{
				Identifier : "09959:114041",
			},
		},
	}

	// Marshal the struct to JSON.
	jsonData, err := json.Marshal(releaseResponse)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate header and status code.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func JFrogHandler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of Response.
	releaseResponse := JFrog{

	Created				: "2025-02-24T05:25:10+00:00",
	Checksums			: map[string]string{
   "sha256":"sammple-sha256",
	},
	DownloadUri			: "sample downlad uri",
	Garbage				: "garbage",
	}

	// Marshal the struct to JSON.
	jsonData, err := json.Marshal(releaseResponse)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate header and status code.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}


func JFrogHandlerItemPropoerties(w http.ResponseWriter, r *http.Request) {
	// Create an instance of Response.
	releaseResponse := ItemProperties{

	JetId				: "sample jet id",
	SealId				: "sample seal id",
	ProjectName 		: "sample project key",
	}

	// Marshal the struct to JSON.
	jsonData, err := json.Marshal(releaseResponse)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate header and status code.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func submitDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate header and status code.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func specificPrefixHandler(prefix string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, prefix) {
			h(w, r)
			return
		}
		http.NotFound(w, r)
	}
}


func specificSuffixHandler(suffix string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.RequestURI, suffix) {
			h(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func specificPrefixAndSuffixHandler(prefix string, suffix string, h http.HandlerFunc, h2 http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.RequestURI, suffix) {
			h2(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI,prefix ) {
			h(w,r)
			return
		}
		http.NotFound(w, r)
	}
}

func main() {
	// r := NewRegexURLMatcher()
		// Register the JSON handler.
		http.HandleFunc("/api/v0.0/releaseready/pipeline/jetId", evidenceCheckHandler)
		// http.HandleFunc("/api/snow/changes(/?[A-Za-z0-9]*)?", snowCheckHandler)

		// http.HandleFunc("/api/storage(/?[A-Za-z0-9]*)?/manifest.json", JFrogHandlerItemPropoerties)
		// http.HandleFunc("/api/storage(/?[A-Za-z0-9]*)?", JFrogHandler)
		http.HandleFunc("/api/snow/changes/", specificPrefixHandler("/api/snow/changes",snowCheckHandler))

		http.HandleFunc("/", specificSuffixHandler("/manifest.json",JFrogHandlerItemPropoerties))
		http.HandleFunc("/api/storage/", specificPrefixAndSuffixHandler("/api/storage","/manifest.json", JFrogHandler,JFrogHandlerItemPropoerties))

		http.HandleFunc("/api/v0.0/deploys/event", submitDeploymentHandler)
		// Start the HTTP server.
		log.Println("Server is running on http://localhost:8093")
		if err := http.ListenAndServe(":8093", nil); err != nil {
			log.Fatalf("Server failed: %s", err)
		}
}


















type RegexURLMatcher struct {
	Patterns map[string]*regexp.Regexp //1
	Handlers map[string]http.HandlerFunc //2
}

func NewRegexURLMatcher() *RegexURLMatcher { //3
	return &RegexURLMatcher{
			Patterns: make(map[string]*regexp.Regexp),
			Handlers: make(map[string]http.HandlerFunc),
	}
}

func (r *RegexURLMatcher) Add(regex string, handler http.HandlerFunc) error { //4        
compiled, err := regexp.Compile(regex) 
	if err != nil {
			return fmt.Errorf("regex string cannot compile with err: %v", compiled)
	}        
	r.Handlers[regex] = handler
	r.Patterns[regex] = compiled        
	return nil
}

func (r *RegexURLMatcher) ServeHTTP(res http.ResponseWriter, req *http.Request) { //5
	toMatchPattern := req.Method + " " + req.URL.Path
	for regexString, handlerFunc := range r.Handlers { //6
			if  r.Patterns[regexString].MatchString(toMatchPattern) {
					handlerFunc(res, req)
					return
			}

	}        
	http.NotFound(res, req)
}