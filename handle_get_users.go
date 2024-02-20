package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

func (cfg *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, u database.User) {
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	res := response{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name.String,
		ApiKey:    u.ApiKey,
	}

	respondWithJSON(w, http.StatusOK, res)
}
