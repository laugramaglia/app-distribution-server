package main

import (
	_ "app-distribution-server-go/docs" // Import the generated docs
	"app-distribution-server-go/internal/application"
	"app-distribution-server-go/internal/infrastructure"
	"app-distribution-server-go/internal/interfaces"
	"log"
	"net/http"
	"regexp"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

// loggingMiddleware logs the incoming requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adds CORS headers to the response.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware: Origin=%s", r.Header.Get("Origin"))
		// Allow requests from any origin. For production, you might want to restrict this.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token")

		// If it's a preflight request, respond with 200 OK
		if r.Method == http.MethodOptions {
			log.Printf("CORS preflight request: Method=%s, Headers=%s", r.Header.Get("Access-Control-Request-Method"), r.Header.Get("Access-Control-Request-Headers"))
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// @title App Distribution API
// @version 1.0
// @description This is a sample server for distributing mobile applications.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	db, err := infrastructure.NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := infrastructure.MigrateDB(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	repo, err := infrastructure.NewPostgresAppRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	service := application.NewAppService(repo)
	handlers := interfaces.NewAppHandlers(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/apps", handlers.AppsHandler)
	mux.HandleFunc("/api/apps/upload", handlers.UploadHandler)
	mux.HandleFunc("/api/apps/", func(w http.ResponseWriter, r *http.Request) {
		downloadRegex := regexp.MustCompile(`/api/apps/([^/]+)/([^/]+)/([^/]+)/download`)
		versionsRegex := regexp.MustCompile(`/api/apps/([^/]+)/versions`)
		latestRegex := regexp.MustCompile(`/api/apps/([^/]+)`)

		if downloadRegex.MatchString(r.URL.Path) {
			handlers.DownloadHandler(w, r)
		} else if versionsRegex.MatchString(r.URL.Path) {
			handlers.GetAllAppVersionsHandler(w, r)
		} else if latestRegex.MatchString(r.URL.Path) {
			handlers.GetLatestAppVersionHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Wrap the mux with the middlewares
	handler := loggingMiddleware(corsMiddleware(mux))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
