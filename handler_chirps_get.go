package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps := []Chirp{}

	chirpsDB, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps from database", err)
		return
	}

	for _, chirpDB := range chirpsDB {
		chirps = append(chirps, Chirp{
			ID:         chirpDB.ID.String(),
			Created_at: chirpDB.CreatedAt,
			Updated_at: chirpDB.UpdatedAt,
			Body:       chirpDB.Body,
			UserID:     chirpDB.UserID.String(),
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
			respondWithError(w, http.StatusBadRequest, "Provided ID doesn't exist", err)
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
