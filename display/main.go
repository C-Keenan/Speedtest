package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"golang.org/x/exp/slices"
)

type Result struct {
	ServerNameandLocation   string
	DownloadMbps float64
	UploadMbps   float64
	PingMs       float64
}

type PageData struct {
	Results []Result
}

func main() {
	const csvFile = "/app/log/ookla_speedtest_log.csv"
	const waitTimeout = 5 * time.Minute
	if err := waitForFile(csvFile, waitTimeout); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/style.css", ServeCss)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func waitForFile(filePath string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			fileInfo, err := os.Stat(filePath)
			if err == nil && fileInfo.Size() > 0 {
				log.Printf("File '%s' found and has content.", filePath)
				return nil
			} else if os.IsNotExist(err) {
				log.Printf("Waiting for file '%s'...", filePath)
			} else if fileInfo.Size() == 0 {
				log.Printf("File '%s' is empty, waiting...", filePath)
			} else {
				log.Printf("Error checking file '%s': %v", filePath, err)
			}
		case <-timeoutChan:
			return fmt.Errorf("timed out waiting for file '%s'", filePath)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	const csvFilePath = "/app/log/ookla_speedtest_log.csv"
	
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Printf("Error opening CSV file at %s: %v", csvFilePath, err)
		http.Error(w, "Could not open CSV file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV data: %v", err)
		http.Error(w, "Could not read CSV data", http.StatusInternalServerError)
		return
	}

	if len(records) <= 1 {
		log.Println("CSV file is empty or contains only headers.")
		pageData := PageData{Results: []Result{}}
		tmpl, _ := template.ParseFiles("template.html")
		tmpl.Execute(w, pageData)
		return
	}

	var results []Result
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 21 {
			log.Printf("Skipping malformed row (expected 21 fields, got %d): %v", len(record), record)
			continue
		}

		serverName := record[0]
		pingStr := record[2]
		downloadStr := record[5]
		uploadStr := record[6]
		
		pingMs, err := strconv.ParseFloat(pingStr, 64)
		if err != nil {
			log.Printf("Error parsing Ping '%s' from record: %v", pingStr, err)
			continue
		}
		
		downloadBps, err := strconv.ParseFloat(downloadStr, 64)
		if err != nil {
			log.Printf("Error parsing Download '%s' from record: %v", downloadBps, err)
			continue
		}
		
		uploadBps, err := strconv.ParseFloat(uploadStr, 64)
		if err != nil {
			log.Printf("Error parsing Upload '%s' from record: %v", uploadBps, err)
			continue
		}
		
		results = append(results, Result{
			ServerNameandLocation:   serverName,
			DownloadMbps: (downloadBps * 8) / 1_000_000,
			UploadMbps:   (uploadBps * 8) / 1_000_000,
			PingMs:       pingMs,
		})
	}
	
	slices.Reverse(results)

	log.Printf("Successfully processed %d data rows for display.", len(results))
	
	pageData := PageData{Results: results}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, pageData)
}

func ServeCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "/static/style.css"))
}