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
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
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
		respondWithError(w, http.StatusInternalServerError, "Problem with feed UUID")
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

	uniqueID, err = uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with feed followUUID")
		return
	}
	feedFollow := database.CreateFeedFollowParams{
		ID:        uniqueID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    uuid.NullUUID{UUID: f.ID, Valid: true},
		UserID:    uuid.NullUUID{UUID: u.ID, Valid: true},
	}
	ff, err := cfg.DB.CreateFeedFollow(r.Context(), feedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem adding feed follow")
	}

	res := response{
		Feed:       f,
		FeedFollow: ff,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Feeds []database.Feed `json:"feeds"`
	}

	res, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Problem retrieving feeds")
	}

	respondWithJSON(w, http.StatusOK, response{Feeds: res})
}
