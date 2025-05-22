package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andras-szesztai/fem_fitness_project/internal/store"
	"github.com/andras-szesztai/fem_fitness_project/internal/tokens"
	"github.com/andras-szesztai/fem_fitness_project/internal/utils"
)

type TokenHandler struct {
	store     store.TokenStore
	userStore store.UserStore
	logger    *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(store store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{store: store, userStore: userStore, logger: logger}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("ERROR: createTokenRequest: %s", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	user, err := th.userStore.GetUserByUsername(req.Username)
	if err != nil {
		th.logger.Printf("ERROR: getUserByUsername: %s", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}

	ok, err := user.PasswordHash.Match(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: matchPassword: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if !ok {
		th.logger.Printf("ERROR: matchPassword: %s", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}

	token, err := th.store.CreateToken(user.ID, 24*time.Hour, tokens.ScopeAuthentication)
	if err != nil {
		th.logger.Printf("ERROR: createToken: %s", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create token"})
		return
	}

	th.logger.Printf("INFO: createToken: %s", token.Plaintext)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}
