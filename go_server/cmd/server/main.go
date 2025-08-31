package main

import (
	"app-distribution-server-go/internal/application"
	"app-distribution-server-go/internal/infrastructure"
	"app-distribution-server-go/internal/interfaces"
	"log"
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

// corsMiddleware adds CORS headers to the response.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin. For production, you might want to restrict this.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token")

		// If it's a preflight request, respond with 200 OK
		if r.Method == http.MethodOptions {
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
	repo, err := infrastructure.NewFileAppRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	service := application.NewAppService(repo)
	handlers := interfaces.NewAppHandlers(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/apps", handlers.AppsHandler)
	mux.HandleFunc("/api/apps/upload", handlers.UploadHandler)
	mux.HandleFunc("/api/apps/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/versions") {
			handlers.GetAllAppVersionsHandler(w, r)
		} else {
			handlers.GetLatestAppVersionHandler(w, r)
		}
	})
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Wrap the mux with the CORS middleware
	corsHandler := corsMiddleware(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
