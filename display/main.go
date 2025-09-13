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
	Timestamp    string
	ServerName   string
	DownloadMbps float64
	UploadMbps   float64
	PingMs       float64
}

type Average struct {
	Period        string
	AvgDownload   float64
	AvgUpload     float64
	AvgPing       float64
	NumDataPoints int
}

type FilterData struct {
	Years  []string
	Months []string
	Days   []string
}

type PageData struct {
	YearlyAverages    []Average
	MonthlyAverages   []Average
	DailyAverages     []Average
	IndividualResults []Result
	FilterData        FilterData
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
	const expectedFields = 22
	
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
		log.Println("CSV file is empty or has only a header.")
		pageData := PageData{}
		tmpl, _ := template.ParseFiles("template.html")
		tmpl.Execute(w, pageData)
		return
	}

	dailyResults := make(map[string][]Result)
	monthlyResults := make(map[string][]Result)
	yearlyResults := make(map[string][]Result)

	var individualResults []Result
	var filteredResults []Result
	
	allYears := make(map[string]bool)
	allMonths := make(map[string]bool)
	allDays := make(map[string]bool)

	selectedYear := r.URL.Query().Get("year")
	selectedMonth := r.URL.Query().Get("month")
	selectedDay := r.URL.Query().Get("day")

	for i, record := range records[1:] {
		if len(record) != expectedFields {
			log.Printf("Skipping malformed row (line %d, expected 22 fields, got %d): %v", i+2, len(record), record)
			continue
		}

		timestampStr := record[21]
		serverName := record[0]
		pingStr := record[2]
		downloadStr := record[5]
		uploadStr := record[6]

		t, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			log.Printf("Error parsing timestamp '%s' on line %d: %v", timestampStr, i+2, err)
			continue
		}

		pingMs, err := strconv.ParseFloat(pingStr, 64)
		if err != nil {
			pingMs = 0.0
		}
		downloadBps, err := strconv.ParseFloat(downloadStr, 64)
		if err != nil {
			downloadBps = 0.0
		}
		uploadBps, err := strconv.ParseFloat(uploadStr, 64)
		if err != nil {
			uploadBps = 0.0
		}

		result := Result{
			Timestamp:    t.Format("2006-01-02 15:04:05 MST"),
			ServerName:   serverName,
			DownloadMbps: (downloadBps * 8) / 1_000_000,
			UploadMbps:   (uploadBps * 8) / 1_000_000,
			PingMs:       pingMs,
		}

		dateKey := t.Format("2006-01-02")
		monthKey := t.Format("2006-01")
		yearKey := t.Format("2006")
		
		dailyResults[dateKey] = append(dailyResults[dateKey], result)
		monthlyResults[monthKey] = append(monthlyResults[monthKey], result)
		yearlyResults[yearKey] = append(yearlyResults[yearKey], result)
		individualResults = append(individualResults, result)

		allYears[t.Format("2006")] = true
		allMonths[t.Format("01")] = true
		allDays[t.Format("02")] = true
	}
	
	for _, res := range individualResults {
		t, err := time.Parse("2006-01-02 15:04:05 MST", res.Timestamp)
		if err != nil {
			continue
		}
		
		matchYear := selectedYear == "" || t.Format("2006") == selectedYear
		matchMonth := selectedMonth == "" || selectedMonth == "any" || t.Format("01") == selectedMonth
		matchDay := selectedDay == "" || selectedDay == "any" || t.Format("02") == selectedDay

		if matchYear && matchMonth && matchDay {
			filteredResults = append(filteredResults, res)
		}
	}

	dailyAverages := calculateAverages(dailyResults)
	monthlyAverages := calculateAverages(monthlyResults)
	yearlyAverages := calculateAverages(yearlyResults)
	
	slices.Reverse(dailyAverages)
	slices.Reverse(monthlyAverages)
	slices.Reverse(yearlyAverages)
	slices.Reverse(filteredResults)

	log.Printf("Successfully processed %d data rows for display.", len(records)-1)

	filterYears := make([]string, 0, len(allYears))
	for y := range allYears { filterYears = append(filterYears, y) }
	slices.Sort(filterYears)

	filterMonths := make([]string, 0, len(allMonths))
	for m := range allMonths { filterMonths = append(filterMonths, m) }
	slices.Sort(filterMonths)

	filterDays := make([]string, 0, len(allDays))
	for d := range allDays { filterDays = append(filterDays, d) }
	slices.Sort(filterDays)
	
	pageData := PageData{
		DailyAverages:    dailyAverages,
		MonthlyAverages:  monthlyAverages,
		YearlyAverages:   yearlyAverages,
		IndividualResults: filteredResults,
		FilterData: FilterData{
			Years:  filterYears,
			Months: filterMonths,
			Days:   filterDays,
		},
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, pageData)
}

func calculateAverages(resultsMap map[string][]Result) []Average {
    var averages []Average
    var keys []string
    for period := range resultsMap {
        keys = append(keys, period)
    }
    slices.Sort(keys)

    for _, period := range keys {
        results := resultsMap[period]
        if len(results) > 0 {
            var totalDownload, totalUpload, totalPing float64
            for _, r := range results {
                totalDownload += r.DownloadMbps
                totalUpload += r.UploadMbps
                totalPing += r.PingMs
            }
            numDataPoints := float64(len(results))
            averages = append(averages, Average{
                Period:        period,
                AvgDownload:   totalDownload / numDataPoints,
                AvgUpload:     totalUpload / numDataPoints,
                AvgPing:       totalPing / numDataPoints,
                NumDataPoints: int(numDataPoints),
            })
        }
    }
	return averages
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "/static/css/style.css"))
}

func serveSetTimedReload(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(".", "/static/js/settimedreload.js")
    http.ServeFile(w, r, filePath)
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Println("Serving file:", filePath)

}

func main() {
	const csvFile = "/app/log/ookla_speedtest_log.csv"
	const waitTimeout = 5 * time.Minute
	if err := waitForFile(csvFile, waitTimeout); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/style.css", serveCss)
	http.HandleFunc("/settimedreload.js", serveSetTimedReload)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
