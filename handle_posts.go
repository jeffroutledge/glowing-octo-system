package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

func (cfg *apiConfig) handlerGetAllPosts(w http.ResponseWriter, r *http.Request, u database.User) {
	type response struct {
		Posts []database.Post
	}

	limitString := chi.URLParam(r, "limit")
	limit := 10
	var err error
	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Problem with limit")
			return
		}
	}

	params := database.GetPostsByUserParams{
		UserID: uuid.NullUUID{UUID: u.ID, Valid: true},
		Limit:  int32(limit),
	}
	posts, err := cfg.DB.GetPostsByUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, response{Posts: posts})
}
