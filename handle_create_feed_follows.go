package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

func (cfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
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
	feedFollow := database.CreateFeedFollowParams{
		ID:        uniqueID,
		FeedID:    uuid.NullUUID{UUID: params.FeedID, Valid: true},
		UserID:    uuid.NullUUID{UUID: u.ID, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ff, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't follow feed")
		return
	}

	res := response{
		ID:        ff.ID,
		FeedID:    ff.FeedID.UUID,
		UserID:    ff.UserID.UUID,
		CreatedAt: ff.CreatedAt,
		UpdatedAt: ff.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, res)
}
