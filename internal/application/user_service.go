package application

import (
	"context"
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) UpdateProfile(ctx context.Context, email, newEmail, newPassword string) error {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	if newEmail != "" {
		user.Email = newEmail
	}

	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("error al encriptar contraseña")
		}
		user.PasswordHash = string(hashedPassword)
	}

	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteProfile(ctx context.Context, email string) error {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("usuario no encontrado")
	}
	return s.repo.Delete(ctx, user.ID)
}
