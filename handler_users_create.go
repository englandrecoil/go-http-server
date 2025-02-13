package main

import (
	"encoding/json"
	"net/http"

	"github.com/englandrecoil/go-http-server/internal/auth"
	"github.com/englandrecoil/go-http-server/internal/database"
	"github.com/lib/pq"
)

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {

	// get request data (email for user)
	userVal := userReqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userVal); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request's parameters", err)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(userVal.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userVal.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				respondWithError(w, http.StatusBadRequest, "This email's already in use", err)
				return
			}
		}
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
