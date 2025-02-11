package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lib/pq"
)

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type userVals struct {
		Email string `json:"email"`
	}
	type userData struct {
		Id         string    `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	// get request data (email for user)
	userVal := userVals{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userVal)

	user, err := cfg.db.CreateUser(r.Context(), userVal.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				respondWithError(w, http.StatusBadRequest, "This email's already in use", err)
				return
			}
		}
	}

	respondWithJSON(w, http.StatusCreated, userData{
		Id:         user.ID.String(),
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})
}
