package application

import (
	"context"
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetProfile recupera la información del usuario por su ID
func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

// UpdateProfile actualiza los datos del usuario (en este caso, su correo)
func (s *UserService) UpdateProfile(ctx context.Context, userID, newEmail string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	user.Email = newEmail
	return s.repo.Update(ctx, user)
}
