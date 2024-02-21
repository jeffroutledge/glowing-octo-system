package main

import (
	"net/http"

	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

func (cfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	ff, err := cfg.DB.GetFeedFollows(r.Context())
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ff)
}
