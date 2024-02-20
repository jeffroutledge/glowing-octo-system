package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	uniqueID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with UUID")
		return
	}
	feed := database.CreateFeedParams{
		ID:        uniqueID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      sql.NullString{String: params.Name, Valid: true},
		Url:       sql.NullString{String: params.Url, Valid: true},
		UserID:    uuid.NullUUID{UUID: u.ID, Valid: true},
	}

	f, err := cfg.DB.CreateFeed(r.Context(), feed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem adding feed")
	}

	res := response{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		Name:      f.Name.String,
		Url:       f.Url.String,
		UserID:    f.UserID.UUID,
	}

	respondWithJSON(w, http.StatusOK, res)
}
