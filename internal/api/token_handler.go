package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zidariu-sabin/femProject/internal/store"
	"github.com/zidariu-sabin/femProject/internal/tokens"
	"github.com/zidariu-sabin/femProject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

// Handler created for uses such as creating authentification token
func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		th.logger.Printf("ERROR: handleCreateToken: %v ", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
	}

	//get user details
	user, err := th.userStore.GetUserByUsername(req.Username)

	if err != nil || user == nil {
		th.logger.Printf("ERROR: GetByUsername: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	//validate password
	passwordDoMatches, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		th.logger.Printf("ERROR: PasswordHash.Matches: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordDoMatches {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := th.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		th.logger.Printf("ERROR: creatingToken: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
