package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/englandrecoil/go-http-server/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID         string    `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	UserID     string    `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	// validate length of chirp
	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// send data to db
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID from string", err)
		return
	}

	createdChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	// send response to client
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:         createdChirp.ID.String(),
		Created_at: createdChirp.CreatedAt,
		Updated_at: createdChirp.UpdatedAt,
		Body:       createdChirp.Body,
		UserID:     createdChirp.UserID.String(),
	})
}

func validateChirp(message string) (string, error) {
	if len(message) > 140 {
		return "", errors.New("chirp is too long")
	}
	return removeProfanity(message), nil

}

func removeProfanity(message string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	originalMsg := strings.Split(message, " ")
	for index, word := range originalMsg {
		loweredWord := strings.ToLower(word)
		if _, ok := profaneWords[loweredWord]; ok {
			originalMsg[index] = "****"
		}
	}
	return strings.Join(originalMsg, " ")

}
