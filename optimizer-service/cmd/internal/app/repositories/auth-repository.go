// Package repositories
package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"optimizer-service/cmd/internal/types"

	"gorm.io/gorm"
)

// AuthRepository is a struct that defines the auth repository
type AuthRepository struct {
	DB *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

// @note: We Can establish any kind of communication we want here,
// with the user service or any other service

// Should return user type or error
func (r *AuthRepository) LoginWithREST(email, password string) (interface{}, error) {
	userServiceUrl := "http://user-service:8080/login"

	loginPayload := types.LoginInput{
		Email:    email,
		Password: password,
	}

	// Converting the payload to bytes, by marshalling it
	payloadBytes, err := json.Marshal(loginPayload)

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Setting up the request
	req, err := http.NewRequest(http.MethodPost, userServiceUrl, bytes.NewBuffer(payloadBytes))

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Setting the headers
	req.Header.Set("Content-Type", "application/json")

	//Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Close the response body
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Invalid credentials")
	}

	var responsePayload types.ResponsePayload

	// read the response body and unmarshal it
	err = json.NewDecoder(resp.Body).Decode(&responsePayload)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return responsePayload.Data, nil
}

// Would validate the token and return a interface and error
// @note: We Can establish any kind of communication we want here,
func (r *AuthRepository) ValidateToken(token string) (interface{}, error) {
	validateUrl := "http://user-service:8080/validate"

	//get the token from the arg
	var requestPayload map[string]string
	requestPayload = make(map[string]string)
	requestPayload["token"] = token

	// Converting the payload to bytes, by marshalling it
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//CReate Post request
	req, err := http.NewRequest(http.MethodPost, validateUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//Set the header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//close the response body
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid token")
	}

	var responsePayload types.ResponsePayload

	err = json.NewDecoder(resp.Body).Decode(&responsePayload)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return responsePayload.Data, nil
}
