package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/englandrecoil/go-http-server/internal/auth"
	"github.com/englandrecoil/go-http-server/internal/database"
	"github.com/google/uuid"
)

type userReqParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUpdateCredentials(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userJWTID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	userVal := userReqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userVal); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request's parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userVal.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash provided password", err)
		return
	}

	userDB, err := cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		HashedPassword: hashedPassword,
		Email:          userVal.Email,
		UpdatedAt:      time.Now().UTC(),
		ID:             userJWTID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user's credentials", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          userDB.ID,
		CreatedAt:   userDB.CreatedAt,
		UpdatedAt:   userDB.UpdatedAt,
		Email:       userDB.Email,
		IsChirpyRed: userDB.IsChirpyRed,
	})
}
