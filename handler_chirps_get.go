package main

import (
	"database/sql"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps := []Chirp{}

	chirpsDB, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps from database", err)
		return
	}

	sortValue := "asc"
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" && (sortParam == "asc" || sortParam == "desc") {
		sortValue = sortParam
	}

	authorIDString := r.URL.Query().Get("author_id")
	authorID := uuid.Nil
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	for _, chirpDB := range chirpsDB {
		if authorID != uuid.Nil && authorID != chirpDB.UserID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:         chirpDB.ID.String(),
			Created_at: chirpDB.CreatedAt,
			Updated_at: chirpDB.UpdatedAt,
			Body:       chirpDB.Body,
			UserID:     chirpDB.UserID.String(),
		})
	}

	if sortValue == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Created_at.Before(chirps[j].Created_at)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Created_at.After(chirps[j].Created_at)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDFromRequst, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID from string", err)
		return
	}

	chirpDB, err := cfg.db.GetChirp(r.Context(), chirpIDFromRequst)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Provided ID doesn't exist", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirp from database", err)
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:         chirpDB.ID.String(),
		Created_at: chirpDB.CreatedAt,
		Updated_at: chirpDB.UpdatedAt,
		Body:       chirpDB.Body,
		UserID:     chirpDB.UserID.String(),
	})

}
