package user

import (
	"fmt"
	"net/http"

	configs "github.com/4lerman/medods_tz/internal/config"
	"github.com/4lerman/medods_tz/internal/service/auth"
	"github.com/4lerman/medods_tz/types"
	"github.com/4lerman/medods_tz/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth", h.handleAuth).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/refresh-token", h.handleRefreshToken).Methods(http.MethodPost)
}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s does not exist", payload.Email))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("token does not match"))
		return
	}

	accessTokenSecret := []byte(configs.Envs.AccessTokenSecret)
	refreshTokenSecret := []byte(configs.Envs.RefreshTokenSecret)
	accessToken, refreshToken, err := auth.CreateJWT(accessTokenSecret, refreshTokenSecret, u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	encryptedToken, err := auth.HashToken(refreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.PutRefreshToken(u.ID, encryptedToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"accessToken": accessToken, "refreshToken": refreshToken})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload types.TokenResponse
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	id := auth.GetUserFromToken(r)
	u, err := h.store.GetUserByID(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with id %v does not exist", id))
		return
	}

	if !auth.CompareTokens(u.RefreshToken, []byte(payload.RefreshToken)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("token does not match"))
		return
	}

	accessTokenSecret := []byte(configs.Envs.AccessTokenSecret)
	refreshTokenSecret := []byte(configs.Envs.RefreshTokenSecret)
	accessToken, refreshToken, err := auth.CreateJWT(accessTokenSecret, refreshTokenSecret, u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	encryptedToken, err := auth.HashToken(refreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.PutRefreshToken(u.ID, encryptedToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"accessToken": accessToken, "refreshToken": refreshToken})

}
