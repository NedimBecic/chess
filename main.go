package main

import (
	"chess/internal/handlers"
	"log"
	"net/http"
	"os"
)

func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedKey := os.Getenv("API_KEY")
		
		if expectedKey != "" {
			apiKey := r.Header.Get("API-KEY")
			if apiKey == "" || apiKey != expectedKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/moves", handlers.HandleGetMoves)
	apiMux.HandleFunc("/move", handlers.HandlePostMove)
	apiMux.HandleFunc("/start", handlers.HandleStartGame)

	fileServer := http.FileServer(http.Dir("./web"))
	
	frontendHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./web/chess.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	mainMux := http.NewServeMux()
	mainMux.Handle("/api/", apiKeyMiddleware(http.StripPrefix("/api", apiMux)))
	mainMux.Handle("/", frontendHandler)

	log.Println("Server starting on :8080")
	log.Println("Frontend: http://localhost:8080/")
	log.Println("API: http://localhost:8080/api/moves (requires API-KEY header)")

	if err := http.ListenAndServe(":8080", mainMux); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}