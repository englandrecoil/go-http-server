package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/englandrecoil/go-http-server/internal/auth"
	"github.com/englandrecoil/go-http-server/internal/database"
)

func (cfg *apiConfig) handlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
	}

	// check if refresh token is valid
	dbRefreshToken, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	diff := dbRefreshToken.ExpiresAt.Sub(time.Now().UTC())
	if diff < 0 {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired", err)
	}
	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is revoked", err)
	}

	// make new access token and send to client
	acessToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't provide new acess token", err)
	}

	type respond struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, respond{
		Token: acessToken,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		UpdatedAt: time.Now().UTC(),
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		Token: refreshToken,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
