package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/englandrecoil/go-http-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get API key", err)
		return
	}
	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong API key provided", err)
		return
	}

	type webhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	params := webhookRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	// find user
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse UUID", err)
		return
	}

	if _, err := cfg.db.GetUserByID(r.Context(), userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't find user", err)
			return
		}
	}

	// upgrade user
	if _, err = cfg.db.UpgradeUser(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
