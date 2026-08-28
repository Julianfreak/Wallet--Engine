package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.PostgresUserRepository
	accountRepo *repository.PostgresAccountRepository // 1. Añadir el repositorio de cuentas
	jwtSecret   string
}

func NewAuthService(userRepo *repository.PostgresUserRepository, accountRepo *repository.PostgresAccountRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, _, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userID := uuid.New().String()
	user := &domain.User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return err
	}

	// 2. Crear automáticamente una cuenta predeterminada para este usuario
	account := &domain.Account{
		ID:      "ACC-" + userID[:8], // ID único basado en su UUID
		Owner:   email,
		Balance: 1000.00, // Saldo inicial de prueba opcional
	}

	return s.accountRepo.Save(ctx, account)
}

// Login valida las credenciales del usuario y retorna un token JWT firmado
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// Compara la contraseña ingresada con el hash guardado en la base de datos
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	fmt.Printf("📦 GENERANDO TOKEN - Usuario ID: %v | Email: '%s'\n", user.ID, user.Email)
	// Generar el Token JWT con una vigencia de 24 horas
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": user.Email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
