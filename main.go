package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jeffroutledge/glowing-octo-system/internal/auth"
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
	dbURL := os.Getenv("CONN")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

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
	v1Router.Post("/users", apiCfg.handlerCreateUsers)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByApiKey))

	app.Mount("/v1", v1Router)

	// srv := &http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: app.HandleFunc,
	// }

	log.Printf("Serving on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, app))
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetApiToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), token)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		handler(w, r, user)
	})
}
