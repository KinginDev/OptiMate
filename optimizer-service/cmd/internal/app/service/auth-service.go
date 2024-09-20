package service

import "optimizer-service/cmd/internal/app/repositories"

type AuthService struct {
	Repo repositories.AuthRepository
}

func NewAuthService(r repositories.AuthRepository) *AuthService {
	return &AuthService{
		Repo: r,
	}
}

func (a *AuthService) Login(email, password string) (string, error) {
	return a.Repo.Login(email, password)
}

func (a *AuthService) ValidateToken(token string) (bool, error) {
	return a.Repo.ValidateToken(token)
}
