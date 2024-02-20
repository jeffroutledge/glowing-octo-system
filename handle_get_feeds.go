package main

import (
	"net/http"

	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

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
