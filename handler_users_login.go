package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/englandrecoil/go-http-server/internal/auth"
	"github.com/englandrecoil/go-http-server/internal/database"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	userVal := userReqParams{}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// get user requests's params
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userVal); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request's parameters", err)
		return
	}

	// search user data by email in db
	userDB, err := cfg.db.FindUserByEmail(r.Context(), userVal.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// authentication part
	err = auth.CheckPasswordHash(userVal.Password, userDB.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// create JWT and refresh token
	accessToken, err := auth.MakeJWT(userDB.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	// inserting refresh token's info in db
	if _, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userDB.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't insert refresh token in db", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          userDB.ID,
			CreatedAt:   userDB.CreatedAt,
			UpdatedAt:   userDB.UpdatedAt,
			Email:       userDB.Email,
			IsChirpyRed: userDB.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
