package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	app := chi.NewRouter()

	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})
	app.Use(corsHandler)

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerError)

	app.Mount("/v1", v1Router)

	// srv := &http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: app.HandleFunc,
	// }

	log.Printf("Serving on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, app))
}
