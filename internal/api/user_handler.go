package api

import (
	"log"
	"net/http"

	"github.com/zidariu-sabin/femProject/internal/store"
)

type registerUserRequest struct {
	Username string,
	
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) handleCreateUser (w http.ResponseWriter, r *http.Request){

}

func (uh *UserHandler) handleGetUserByUsername (w http.ResponseWriter, r *http.Request){

}

func (uh *UserHandler) handleSearch (w http.ResponseWriter, r *http.Request){
	
}