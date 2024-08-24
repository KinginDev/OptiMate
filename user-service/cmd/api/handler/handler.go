package handler

import (
	"net/http"
	"time"
	"user-service/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Handler struct to hold the db instance
type Handler struct {
	DB *gorm.DB
}

// NewHandler function to initialize the handler with the given DB instance
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

type JsonResponse struct {
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
	response := &JsonResponse{
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
// @Success 201 {object} JsonResponse "User created successfully"
// @Failure 400 {object} JsonResponse "Invalid request payload"
// @Failure 400 {object} JsonResponse "Failed to create user"
// @Failure 409 {object} JsonResponse "Email already exists"
// @Router /register [post]
func (h *Handler) Register(c echo.Context) error {
	u := new(models.User)

	/**
		The Bind function helps to map the incoming
		JSON request body to the corresponding struct
	**/
	if err := c.Bind(u); err != nil {
		return err
	}

	//validate input fields
	if err := c.Validate(u); err != nil {
		errResponsePayload := &JsonResponse{
			Data:    "",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		return c.JSON(errResponsePayload.Status, errResponsePayload)
	}

	//check if email already exists
	var count int64
	h.DB.Where("email = ?", u.Email).Find(&models.User{}).Count(&count)
	if count > 0 {
		errResponsePayload := &JsonResponse{
			Data:    "",
			Message: "Email already exists",
			Status:  http.StatusConflict,
		}
		return c.JSON(errResponsePayload.Status, errResponsePayload)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.ID = uuid.New().String()

	if err := h.DB.Create(&u).Error; err != nil {
		errResponsePayload := &JsonResponse{
			Data:    "",
			Message: "Failed to create user",
			Status:  http.StatusBadRequest,
		}
		return c.JSON(errResponsePayload.Status, errResponsePayload)
	}

	type UserJsonResponse struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	userJsonResponse := &UserJsonResponse{
		ID:    u.ID,
		Email: u.Email,
	}

	responsePayload := &JsonResponse{
		Data:    userJsonResponse,
		Message: "User created successfully",
		Status:  http.StatusCreated,
	}
	return c.JSON(responsePayload.Status, responsePayload)

}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Accept json
// @Produce json
// @Success 200 {object} JsonResponse "Login successful"
// @Failure 400 {object} JsonResponse "Invalid request payload"
// @Failure 401 {object} JsonResponse "Invalid password"
// @Failure 404 {object} JsonResponse "User not found"
func (h *Handler) Login(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		errorResponseJson := &JsonResponse{
			Data:    "",
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		return c.JSON(errorResponseJson.Status, errorResponseJson)
	}

	//validate input fields
	if err := c.Validate(u); err != nil {
		errResponsePayload := &JsonResponse{
			Data:    "",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		return c.JSON(errResponsePayload.Status, errResponsePayload)
	}

	var user models.User

	//Query the database to find the user with the email
	if err := h.DB.Where("email = ?", u.Email).First(&user).Error; err != nil {
		errorResponseJson := &JsonResponse{
			Data:    "",
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
		return c.JSON(errorResponseJson.Status, errorResponseJson)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),
		[]byte(u.Password)); err != nil {

		errorResponseJson := &JsonResponse{
			Data:    "",
			Message: "Invalid password",
			Status:  http.StatusUnauthorized,
		}
		return c.JSON(errorResponseJson.Status, errorResponseJson)

	}

	//Generate new JWT token fork user
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		errorResponseJson := &JsonResponse{
			Data:    "",
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
		return c.JSON(errorResponseJson.Status, errorResponseJson)
	}

	personalToken := &models.PersonalToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     t,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	if err := h.DB.Create(personalToken).Error; err != nil {
		errorResponseJson := &JsonResponse{
			Data:    "",
			Message: "Failed to create token",
			Status:  http.StatusInternalServerError,
		}
		return c.JSON(errorResponseJson.Status, errorResponseJson)
	}

	responsePayload := &JsonResponse{
		Data:    personalToken,
		Message: "Login successful",
		Status:  http.StatusOK,
	}

	return c.JSON(responsePayload.Status, responsePayload)
}
