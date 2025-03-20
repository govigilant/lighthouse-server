package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sync"
	"bytes"
)

type Request struct {
	Website     string `json:"website"`
	CallbackURL string `json:"callback_url"`
}

var mu sync.Mutex
var processing bool

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if processing {
		mu.Unlock()
		http.Error(w, "Server is busy", http.StatusServiceUnavailable)
		return
	}
	processing = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		processing = false
		mu.Unlock()
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	go runLighthouse(req)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Processing started"))
}

func runLighthouse(req Request) {
	cmd := exec.Command("lighthouse", req.Website, "--quiet", "--chrome-flags=--headless --no-sandbox --disable-dev-shm-usage --disable-gpu", "--output=json", "--output-path=report.json")
	fmt.Println("Running Lighthouse for", req.Website)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running Lighthouse:", err)
		return
	}

	report, err := ioutil.ReadFile("report.json")
	if err != nil {
		fmt.Println("Error reading Lighthouse report:", err)
		return
	}

	fmt.Println("Sending report to", req.CallbackURL)
	resp, err := http.Post(req.CallbackURL, "application/json", ioutil.NopCloser(bytes.NewReader(report)))

	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func main() {
	http.HandleFunc("/lighthouse", handler)
	fmt.Println("Server started on :80")
	http.ListenAndServe(":80", nil)
}
