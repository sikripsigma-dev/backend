package service

import (
	"Skripsigma-BE/internal/constants"
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/util"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo       repository.UserRepository
	authTokenRepo  repository.AuthTokenRepository
	chatService    *ChatService
	// mailer         *util.Mailer
}

func NewAuthService(userRepo repository.UserRepository, authTokenRepo repository.AuthTokenRepository, chatService *ChatService) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		authTokenRepo: authTokenRepo,
		chatService:   chatService,
		// mailer:        mailer,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*models.User, error) {
	// 1. Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("Email already registered")
	}

	// 2. Cari dosen pembimbing berdasarkan email
	supervisor, err := s.userRepo.GetByEmail(req.SupervisorEmail)
	if err != nil {
		return nil, fmt.Errorf("Dosen pembimbing dengan email tersebut tidak ditemukan")
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password")
	}

	// 4. Simpan user
	user := models.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		// default role mahasiswa 3
		RoleId: 3,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, fmt.Errorf("Failed to register user")
	}

	// 5. Simpan ke tabel user_student dengan supervisor_id
	student := &models.StudentUser{
		UserID:       user.Id,
		SupervisorID: supervisor.Id,
		UniversityID: req.UniversityID,
		Jurusan:      req.Jurusan,
		Gpa:          req.Gpa,
		Nim:          req.Nim,
	}

	if err := s.userRepo.CreateStudent(student); err != nil {
		return nil, fmt.Errorf("Gagal menyimpan data mahasiswa")
	}

	// 6. Generate token verifikasi email
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	tokenData := &models.AuthToken{
		UserID:    user.Id,
		Token:     token,
		Type:      "email_verification",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	_ = s.authTokenRepo.Create(tokenData) // Boleh abaikan error di awal

	// 7. Kirim email verifikasi
	go func() {
		mailer := util.NewMailer()
		// verifyURL := fmt.Sprintf("https://skripsigma.com/verify?token=%s", token)
		frontendURL := os.Getenv("FRONTEND_URL") // Contoh: http://localhost:3000
		verifyURL := fmt.Sprintf("%s/verify?token=%s", frontendURL, token)
		subject := "Verifikasi Email Anda - Skripsigma"
		body := fmt.Sprintf("Halo %s,\n\nSilakan klik link berikut untuk verifikasi email Anda:\n\n%s\n\nLink berlaku 24 jam.\n\nTerima kasih.", user.Name, verifyURL)
		_ = mailer.Send(user.Email, subject, body)
	}()

	return &user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("Invalid email or password")
	}

	// (Optional) bisa tambahkan pengecekan apakah sudah verifikasi
	if !user.IsVerified { return "", fmt.Errorf("Email belum diverifikasi") }

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
	user, err := s.userRepo.GetWithRelationsByID(userId)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return user, nil
}

// func (s *AuthService) VerifyEmailToken(token string) (*models.AuthToken, error) {
// 	authToken, err := s.authTokenRepo.GetByToken(token)
// 	if err != nil {
// 		return nil, fmt.Errorf("Token tidak valid")
// 	}

// 	if time.Now().After(authToken.ExpiresAt) {
// 		return nil, fmt.Errorf("Token sudah kedaluwarsa")
// 	}

// 	// Update user menjadi terverifikasi
// 	user, err := s.userRepo.GetByID(authToken.UserID)
// 	if err != nil {
// 		return nil, fmt.Errorf("User tidak ditemukan")
// 	}

// 	user.IsVerified = true
// 	if err := s.userRepo.Update(user); err != nil {
// 		return nil, fmt.Errorf("Gagal memverifikasi email")
// 	}

// 	_ = s.authTokenRepo.DeleteByID(authToken.ID)

// 	return authToken, nil
// }

func (s *AuthService) VerifyEmailToken(token string) (*models.AuthToken, error) {
	authToken, err := s.authTokenRepo.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("Token tidak valid")
	}

	if time.Now().After(authToken.ExpiresAt) {
		return nil, fmt.Errorf("Token sudah kedaluwarsa")
	}

	// Update user menjadi terverifikasi
	user, err := s.userRepo.GetByID(authToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("User tidak ditemukan")
	}

	user.IsVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("Gagal memverifikasi email")
	}

	// 🔹 Auto create chat room jika mahasiswa
	if user.RoleId == constants.RoleStudent {
		studentDetail, err := s.userRepo.GetStudentDetailByID(user.Id)
		if err == nil && studentDetail.Student != nil && studentDetail.Student.SupervisorID != "" {
			_, _ = s.chatService.CreateOrGetChatRoom(user.Id, &studentDetail.Student.SupervisorID, nil)
		}
	}

	_ = s.authTokenRepo.DeleteByID(authToken.ID)

	return authToken, nil
}


