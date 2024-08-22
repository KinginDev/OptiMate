package handler

import (
	"net/http"
	"time"
	"user-service/data"

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

func (h *Handler) Index(c echo.Context) error {
	response := &JsonResponse{
		Data:    "Welcome to the User Service!",
		Message: "Service is running.",
		Status:  http.StatusOK,
	}
	return c.JSON(response.Status, response)
}

func (h *Handler) Register(c echo.Context) error {
	u := new(data.User)

	/**
		The Bind function helps to map the incoming
		JSON request body to the corresponding struct
	**/
	if err := c.Bind(u); err != nil {
		return err
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
			Status:  http.StatusInternalServerError,
		}
		return c.JSON(errResponsePayload.Status, errResponsePayload)
	}

	responsePayload := &JsonResponse{
		Data:    u,
		Message: "User created successfully",
		Status:  http.StatusCreated,
	}
	return c.JSON(responsePayload.Status, responsePayload)

}

func (h *Handler) Login(c echo.Context) error {
	u := new(data.User)
	if err := c.Bind(u); err != nil {
		errorResponseJson := &JsonResponse{
			Data:    "",
			Message: "Invalid request payload",
			Status:  http.StatusBadRequest,
		}
		return c.JSON(errorResponseJson.Status, errorResponseJson)
	}

	var user data.User

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

	personalToken := &data.PersonalToken{
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
