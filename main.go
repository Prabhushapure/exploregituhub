package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const githubAPIURL = "https://api.github.com"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/repositories/{owner}/{repo}", getRepositoryInfo).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, router)
}

func getRepositoryInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	owner := params["owner"]
	repo := params["repo"]

	url := fmt.Sprintf("%s/repos/%s/%s", githubAPIURL, owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Repository not found", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
