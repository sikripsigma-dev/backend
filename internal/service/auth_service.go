package service

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*models.User, error) {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil { // Pastikan user ditemukan
		return nil, fmt.Errorf("Email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password")
	}

	user := models.User{
		Nim:     req.Nim,
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Password: string(hashedPassword),
	}

	// Simpan user ke DB
	if err := s.userRepo.Create(&user); err != nil {
		return nil, fmt.Errorf("Failed to register user")
	}

	return &user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (string, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("Invalid email or password")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("Invalid email or password")
	}

	// Buat token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.RegisteredClaims{
		Issuer:    "Skripsigma-API",
		Subject:   user.Id,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("Failed to generate token")
	}

	return t, nil
}

func (s *AuthService) GetUserByToken(tokenString string) (*models.User, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	userId := claims.Subject

	user, err := s.userRepo.GetWithCompanyByID(userId)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return user, nil
}
