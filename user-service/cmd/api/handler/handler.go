// Package handler contains the route functions for the user service
package handler

import (
	"net/http"
	"time"
	"user-service/models"
	"user-service/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Handler struct to hold the db instance
type Handler struct {
	DB    *gorm.DB
	Utils *utils.Utils
}

type UserJSONResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// NewHandler function to initialize the handler with the given DB instance
func NewHandler(db *gorm.DB, utils *utils.Utils) *Handler {
	return &Handler{
		DB:    db,
		Utils: utils,
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
	u := new(models.User)

	/**
		The Bind function helps to map the incoming
		JSON request body to the corresponding struct
	**/
	if err := c.Bind(u); err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// validate input fields
	if err := c.Validate(u); err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// check if email already exists
	var count int64
	h.DB.Where("email = ?", u.Email).Find(&models.User{}).Count(&count)
	if count > 0 {
		return h.Utils.WriteErrorResponse(c, http.StatusConflict, "Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Failed to create user")
	}

	u.Password = string(hashedPassword)
	u.ID = uuid.New().String()

	// Create the user
	err = u.CreateUser(h.DB, u.Email, u.Password)
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Failed to create user")
	}

	UserJSONResponse := &UserJSONResponse{
		ID:    u.ID,
		Email: u.Email,
	}

	return h.Utils.WriteSuccessResponse(c, http.StatusCreated, "User created successfully", UserJSONResponse)

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
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// validate input fields
	if err := c.Validate(u); err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Query the database to find the user with the email
	user, err := u.GetUserByEmail(h.DB, u.Email)
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusNotFound, "User not found")
	}

	// compare the password
	if err := user.ComparePassword(u.Password); !err {
		return h.Utils.WriteErrorResponse(c, http.StatusUnauthorized, "Invalid password")
	}

	// Generate new JWT token for user
	t, err := h.Utils.GenerateJWTToken(user.ID)
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}

	personalToken := &models.PersonalToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     t,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	if err := h.DB.Create(personalToken).Error; err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to create token")
	}

	return h.Utils.WriteSuccessResponse(c, http.StatusOK, "Login successful", map[string]string{
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
func (h *Handler) GetUserTokens(c echo.Context) error {
	userID, ok := c.Get("userID").(string)

	if !ok {
		return h.Utils.WriteErrorResponse(c, http.StatusBadRequest, "User not found")
	}

	var u models.User

	user, err := u.GetUserByID(h.DB, userID)
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusNotFound, "User not found")
	}

	tokens, err := models.GetTokensByUserID(h.DB, user.ID)
	if err != nil {
		return h.Utils.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tokens")
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
		"tokens": tokens[len(tokens)-1],
	}

	return h.Utils.WriteSuccessResponse(c, http.StatusOK, "User tokens retrieved successfully", response)
}
