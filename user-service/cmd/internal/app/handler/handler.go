// Package handler contains the route functions for the user service
package handler

import (
	"net/http"
	"time"
	"user-service/cmd/internal/models"
	"user-service/cmd/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler struct to hold the db instance
type Handler struct {
	Container *types.AppContainer
}

type UserJSONResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// NewHandler function to initialize the handler with the given DB instance
func NewHandler(container *types.AppContainer) *Handler {
	return &Handler{
		Container: container,
	}
}

type JSONResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

// Index godoc
// @Summary User service is running
// @Description User service is running
// @Success 200 {object} nil "success"
// @Failure 404 {object} nil "Not Found"
// @Router / [get]
func (h *Handler) Index(c echo.Context) error {
	response := &JSONResponse{
		Data:    "Welcome to the User Service!",
		Message: "Service is running.",
		Status:  http.StatusOK,
	}
	return c.JSON(response.Status, response)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept json
// @Produce json
// @Success 201 {object} JSONResponse "User created successfully"
// @Failure 400 {object} JSONResponse "Invalid request payload"
// @Failure 400 {object} JSONResponse "Failed to create user"
// @Failure 409 {object} JSONResponse "Email already exists"
// @Router /register [post]
func (h *Handler) Register(c echo.Context) error {
	var input models.RegisterInput
	/**
		The Bind function helps to map the incoming
		JSON request body to the corresponding struct
	**/
	if err := c.Bind(&input); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// validate input fields
	if err := c.Validate(&input); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	existingUser, err := h.Container.UserService.Repo.GetUserByEmail(input.Email)

	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if existingUser != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Email already exists")
	}

	user, err := h.Container.UserService.RegisterUser(&input)

	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Generate new JWT token for user
	t, err := h.Container.JWTService.GenerateJWTToken(user.ID)
	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}

	personalToken := &models.PersonalToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     t,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	if err := h.Container.JWTService.StoreToken(personalToken); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}

	UserJSONResponse := &UserJSONResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusCreated, "User created successfully", UserJSONResponse)

}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Accept json
// @Produce json
// @Success 200 {object} JSONResponse "Login successful"
// @Failure 400 {object} JSONResponse "Invalid request payload"
// @Failure 401 {object} JSONResponse "Invalid password"
// @Failure 404 {object} JSONResponse "User not found"
func (h *Handler) Login(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// validate input fields
	if err := c.Validate(u); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.Container.UserService.AuthenticateUser(u.Email, u.Password)
	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusUnauthorized, "Invalid Credentials")
	}

	// Generate new JWT token for user
	t, err := h.Container.JWTService.GenerateJWTToken(user.ID)
	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}

	personalToken := &models.PersonalToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     t,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	if err := h.Container.JWTService.StoreToken(personalToken); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "Login successful", map[string]string{
		"email": user.Email,
		"token": t,
	})

}

// GetUserTokens godoc
// @Summary  get toksn
// @Description gets all the authorization token belongin to a user
// @Accept json
// @Produce json
// @Success 200 (object) JSONResponse "User tokens retrieved successfully"
// @Faliure 400 (object) JSONResponse "User not found"
// @Faliure 404 (object) JSONResponse "User not found"
// @Faliure 500 (object) JSONResponse "Failed to fetch tokens"
func (h *Handler) GetUserJWTTokens(c echo.Context) error {

	// Get the user id from the middleware
	userID, ok := c.Get("userID").(string)

	if !ok {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, "User not found")
	}

	user, err := h.Container.UserService.Repo.GetUserByID(h.Container.DB, userID)
	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusNotFound, "User not found")
	}

	tokens, err := h.Container.UserService.GetUserTokens(user.ID)
	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tokens")
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
		"tokens": tokens[len(tokens)-1],
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "User tokens retrieved successfully", response)
}
