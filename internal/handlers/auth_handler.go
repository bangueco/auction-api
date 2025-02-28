package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/bangueco/auction-api/internal/handlers/helper"
	"github.com/bangueco/auction-api/internal/lib"
	"github.com/bangueco/auction-api/internal/models"
	"github.com/bangueco/auction-api/internal/services"
)

type AuthHandler struct {
	UserService *services.UserService
}

func NewAuthHandler(UserService *services.UserService) *AuthHandler {
	return &AuthHandler{UserService}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := helper.DecodeRequestBody(r, &user)

	if err != nil {
		helper.WriteResponseMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	errorMessages := lib.ValidateStruct(&user)

	if errorMessages != nil {
		helper.WriteResponse(w, errorMessages, http.StatusBadRequest)
		return
	}

	_, err = a.UserService.GetUserByUsername(user.Username)

	if !errors.Is(err, services.ErrUserNotFound) {
		helper.WriteResponseMessage(w, "Username is already taken", http.StatusConflict)
		return
	}

	newUser, err := a.UserService.CreateUser(user)

	if err != nil {
		helper.WriteResponseMessage(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	helper.WriteResponse(w, newUser, http.StatusCreated)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := helper.DecodeRequestBody(r, &user)

	if err != nil {
		log.Printf("Cannot decode json body: %s", err)
		helper.WriteResponseMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	errorMessages := lib.ValidateStruct(&user)

	if errorMessages != nil {
		helper.WriteResponse(w, errorMessages, http.StatusBadRequest)
		return
	}

	existingUser, err := a.UserService.GetUserByUsername(user.Username)

	if errors.Is(err, services.ErrUserNotFound) {
		helper.WriteResponseMessage(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	match := lib.ComparePassword(user.Password, existingUser.Password)

	if !match {
		helper.WriteResponseMessage(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	token, err := lib.GenerateToken(user.ID)

	if err != nil {
		helper.WriteResponseMessage(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"token": token,
	}

	helper.WriteResponse(w, data, http.StatusOK)
}
