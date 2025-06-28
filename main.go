// @title URL Shortener with Insights API
// @version 1.0
// @description A GoFr-based URL shortening service with analytics
// @host localhost:8080
// @BasePath /
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/internal/utils"

	//_ "url-shortener/docs"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	shortenedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortened_total",
			Help: "Total number of short links created",
		},
	)
	redirectCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_redirect_total",
			Help: "Total number of redirects",
		},
	)

	db *sql.DB
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Connects to PostgreSQL
	var err error
	db, err = repository.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create 'urls' and 'visits' tables if they don't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		id TEXT PRIMARY KEY,
		original TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		custom_code TEXT,
		domain TEXT,
		visibility TEXT,
		created_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS visits (
		id SERIAL PRIMARY KEY,
		url_id TEXT REFERENCES urls(id) ON DELETE CASCADE,
		timestamp TEXT,
		ip_address TEXT,
		country TEXT,
		browser TEXT,
		device TEXT
	);
`)

	if err != nil {
		panic("❌ Failed to create tables: " + err.Error())
	}

	fmt.Println("✅ Ensured tables 'urls' and 'visits' exist")
	fmt.Println("✅ Connected to PostgreSQL")

	// Set up HTTP handlers
	// 	http.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
	// 		httpSwagger.Handler(
	// 			httpSwagger.URL("http://localhost:8080/docs/swagger.json"), // ✅ Tell Swagger where to find your spec
	// 		)(w, r)
	// 	})
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/analytics/", analyticsHandler)
	http.HandleFunc("/all", getAllLinksHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	fmt.Println("Server running on :8080")
	fmt.Println(`Example: curl.exe -X POST http://localhost:8080/shorten -H "Content-Type: application/x-www-form-urlencoded" -d "url=https://google.com"`)
	http.ListenAndServe(":8080", nil)
}

func getClientIP(r *http.Request) string {
	// Check for forwarded headers first (if behind proxy)
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// @Summary Get analytics for a short link
// @Description Returns visit analytics for a short code
// @Tags Analytics
// @Produce text/plain
// @Param code path string true "Short code"
// @Success 200 {string} string "Visit analytics"
// @Failure 404 {string} string "Short code not found"
// @Router /analytics/{code} [get]
func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/analytics/")
	if code == "" {
		http.Error(w, "Analytics code required", http.StatusBadRequest)
		return
	}

	// Get URL by short code
	url, err := repository.GetURLByCode(db, code)

	if url.Visibility == "private" {
		http.Error(w, "Analytics not available for private URLs", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	// Get visits for that URL
	visits, err := repository.GetVisitsByURLID(db, url.ID)
	if err != nil {
		http.Error(w, "Could not fetch visits", http.StatusInternalServerError)
		return
	}

	// Return analytics
	if len(visits) == 0 {
		fmt.Fprintln(w, "No visits yet.")
		return
	}

	for i, v := range visits {
		fmt.Fprintf(w, "Visit %d:\n", i+1)
		fmt.Fprintf(w, "  IP        : %s\n", v.IPAddress)
		fmt.Fprintf(w, "  Country   : %s\n", v.Country)
		fmt.Fprintf(w, "  Timestamp : %s\n", v.Timestamp)
		fmt.Fprintf(w, "  Browser   : %s\n", v.Browser)
		fmt.Fprintf(w, "  Device    : %s\n\n", v.Device)
	}

}

// @Summary Create a shortened URL
// @Description Generates a short link for the provided original URL
// @Tags Shortener
// @Accept  application/x-www-form-urlencoded
// @Produce text/plain
// @Param url formData string true "Original URL to shorten"
// @Param code formData string false "Custom short code"
// @Param visibility formData string false "public or private"
// @Success 200 {string} string "Short URL and Analytics URL"
// @Failure 400 {string} string "Bad request"
// @Router /shorten [post]
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	original := r.PostFormValue("url")
	if original == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// Read optional custom short code
	requestedCode := r.FormValue("code")
	var shortCode string

	fmt.Println("Custom code from request:", requestedCode)

	if requestedCode != "" {
		// Check if code already exists
		_, err := repository.GetURLByCode(db, requestedCode)
		if err == nil {
			http.Error(w, "Custom short code already in use", http.StatusBadRequest)
			return
		}
		shortCode = requestedCode
	} else {
		shortCode = generateUniqueShortCode(6)
	}

	visibility := r.PostFormValue("visibility")
	if visibility != "private" {
		visibility = "public" // default fallback
	}

	url := model.URL{
		ID:         uuid.New().String(), // ← this is a unique ID
		Original:   original,
		ShortCode:  shortCode,
		Visibility: visibility,
		CreatedAt:  model.Now(),
	}

	// Save to DB
	err = repository.SaveURL(db, url)
	if err != nil {
		http.Error(w, "Failed to save to database", http.StatusInternalServerError)
		fmt.Println("DB Error:", err)
		return
	}
	shortenedCounter.Inc()

	// mutex.Lock()
	// urlMap[shortCode] = original
	// analytics[shortCode] = []map[string]interface{}{} // ✅ Init analytics to avoid Not Found
	// mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url":     "http://localhost:8080/r/" + shortCode,
		"analytics_url": "http://localhost:8080/analytics/" + shortCode,
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/r/")
	if code == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Fetch from DB
	url, err := repository.GetURLByCode(db, code)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Track visit (we’ll enhance this later)
	// Get IP (client or fallback)
	ip := getClientIP(r)
	if ip == "::1" || ip == "127.0.0.1" {
		ip = "103.48.198.141" // local dev IP override
	}

	// Parse user-agent
	ua := r.UserAgent()
	browser, os, device := utils.ParseUserAgent(ua)

	// Get location
	location, _ := utils.GetLocation(ip) // ignore error fallback

	visit := model.Visit{
		URLID:     url.ID,
		Timestamp: time.Now(),
		IPAddress: ip,
		Country:   location.Country,
		//City:      location.City,
		Browser: browser,
		OS:      os,
		Device:  device,
	}

	err = repository.SaveVisit(db, visit)
	if err != nil {
		fmt.Println("❌ Failed to log visit:", err)
	} else {
		fmt.Println("✅ Visit logged")
	}

	redirectCounter.Inc()

	// Step 3: Redirect
	http.Redirect(w, r, url.Original, http.StatusFound)
}

func generateShortCode(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func generateUniqueShortCode(length int) string {
	for {
		code := generateShortCode(length)
		exists, err := repository.ShortCodeExists(db, code)
		if err != nil {
			fmt.Println("DB error while checking short code:", err)
			continue
		}
		if !exists {
			return code
		}
	}
}

// @Summary Get all short links
// @Description Lists all short links in the database
// @Tags Links
// @Produce text/plain
// @Success 200 {string} string "List of short links"
// @Router /all [get]
func getAllLinksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	links, err := repository.GetAllLinks(db)
	if err != nil {
		http.Error(w, "Failed to fetch links", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

// @Summary Update a short link
// @Description Update original URL or visibility
// @Tags Links
// @Accept application/x-www-form-urlencoded
// @Param code path string true "Short code"
// @Param url formData string false "New original URL"
// @Param visibility formData string false "New visibility"
// @Success 200 {string} string "Updated short link"
// @Failure 404 {string} string "Short code not found"
// @Router /update/{code} [put]
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/update/")
	if code == "" {
		http.Error(w, "Short code required", http.StatusBadRequest)
		return
	}

	// Decode JSON from request body
	var req struct {
		LongURL    string `json:"long_url"`
		Visibility string `json:"visibility"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Fetch original short URL from DB
	url, err := repository.GetURLByCode(db, code)
	if err != nil {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	// Apply updates if provided
	if req.LongURL != "" {
		url.Original = req.LongURL
	}
	if req.Visibility != "" {
		url.Visibility = req.Visibility
	}

	// Save updates
	err = repository.UpdateURL(db, url)
	if err != nil {
		http.Error(w, "Failed to update URL", http.StatusInternalServerError)
		return
	}

	// Respond success
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "✅ Updated short link '%s'", code)
}

// @Summary Delete a short link
// @Description Permanently deletes a short code
// @Tags Links
// @Param code path string true "Short code"
// @Success 200 {string} string "Deleted short link"
// @Failure 404 {string} string "Short code not found"
// @Router /delete/{code} [delete]
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/delete/")
	if code == "" {
		http.Error(w, "Short code required", http.StatusBadRequest)
		return
	}

	url, err := repository.GetURLByCode(db, code)
	if err != nil {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	fmt.Println("⛏️ Attempting to delete short code:", code, "| ID:", url.ID)

	err = repository.DeleteURLByCode(db, code)
	if err != nil {
		fmt.Println("❌ DB Delete error:", err) // <--- Add this line
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "🗑️ Deleted short link '%s'", code)
}

// @Summary Health check
// @Description Returns service health status
// @Tags Monitoring
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Database unreachable"
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// optional DB ping check
	err := db.Ping()
	if err != nil {
		http.Error(w, `{"status":"db unreachable"}`, http.StatusInternalServerError)
		return
	}

	// success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
