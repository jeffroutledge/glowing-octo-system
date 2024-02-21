package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

func (cfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	type response struct{}

	feedFollowIdString := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem with feed follow id")
		return
	}

	_, err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusNoContent, response{})
}
