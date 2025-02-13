package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/englandrecoil/go-http-server/internal/auth"
	"github.com/englandrecoil/go-http-server/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userReqParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

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
	}

	// inserting refresh token's info in db
	cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userDB.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	})

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userDB.ID,
			CreatedAt: userDB.CreatedAt,
			UpdatedAt: userDB.UpdatedAt,
			Email:     userDB.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
