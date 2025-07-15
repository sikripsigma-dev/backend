package handler

import (
	"log"
	"net/http"
	"time"

	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCreateHandler struct {
	UserCreateService service.UserCreateService
	UserCreateLogSvc  service.UserCreateLogService
}

func NewUserCreateHandler(userSvc service.UserCreateService, logSvc service.UserCreateLogService) *UserCreateHandler {
	return &UserCreateHandler{
		UserCreateService: userSvc,
		UserCreateLogSvc:  logSvc,
	}
}

func (h *UserCreateHandler) CreateUser(c *fiber.Ctx) error {
	type Request struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		RoleId       uint   `json:"role_id"`
		Password     string `json:"password"`
		Nidn         string `json:"nidn"`
		UniversityID string `json:"university_id"`
		CompanyID    string `json:"company_id"`
		Division     string `json:"division"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request",
			"details": err.Error(),
		})
	}

	// Get creator user from context
	creator, ok := c.Locals("user").(*models.User)
	if !ok || creator == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized or user context missing",
		})
	}

	// Generate UUID once and use it consistently
	userID := uuid.New().String()
	log.Printf("Creating new user with ID: %s", userID)

	// Validate role-specific fields
	if req.RoleId == 4 && (req.UniversityID == "" || req.Nidn == "") {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "University ID and NIDN are required for supervisor role",
		})
	}
	if req.RoleId == 2 && (req.CompanyID == "" || req.Division == "") {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Company ID and division are required for company user role",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password",
		})
	}

	// Prepare user model - ID will be preserved because we set it
	newUser := &models.User{
		Id:       userID, // This will NOT be overwritten due to our updated BeforeCreate hook
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleId:   req.RoleId,
		Password: string(hashedPassword),
	}

	// Prepare role-specific models
	var supervisor *models.SupervisorUser
	var companyUser *models.CompanyUser

	switch req.RoleId {
	case 4: // Supervisor
		supervisor = &models.SupervisorUser{
			UserID:       userID,
			UniversityID: req.UniversityID,
			Nidn:         req.Nidn,
		}
		log.Printf("Creating supervisor record for user: %s", userID)
	case 2: // Company user
		companyUser = &models.CompanyUser{
			UserID:    userID,
			CompanyID: req.CompanyID,
			Division:  req.Division,
		}
		log.Printf("Creating company user record for user: %s", userID)
	}

	// Create user and related records
	if err := h.UserCreateService.CreateUserFull(newUser, supervisor, companyUser); err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
	}

	// Create audit log
	logEntry := models.UserCreateLog{
		ID:        uuid.New().String(),
		CreatedBy: creator.Id,
		UserID:    userID,
		RoleID:    req.RoleId,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.UserCreateLogSvc.CreateLog(&logEntry); err != nil {
		log.Printf("Failed to create user creation log: %v", err)
		// Don't fail the request if log creation fails
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user_id": userID,
		"email":   req.Email,
	})
}


func (h *UserCreateHandler) GetCreationLogs(c *fiber.Ctx) error {
    logs, err := h.UserCreateLogSvc.GetLogs()
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch logs",
            "details": err.Error(),
        })
    }
    return c.JSON(logs)
}

func (h *UserCreateHandler) GetMyCreatedUsers(c *fiber.Ctx) error {
    user, ok := c.Locals("user").(*models.User)
    if !ok || user == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unauthorized",
        })
    }

    logs, err := h.UserCreateLogSvc.GetLogsByCreator(user.Id)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch your created users",
            "details": err.Error(),
        })
    }
    return c.JSON(logs)
}

func (h *UserCreateHandler) GetUserCreationHistory(c *fiber.Ctx) error {
    userID := c.Params("id")
    logs, err := h.UserCreateLogSvc.GetLogsByUser(userID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch user creation history",
            "details": err.Error(),
        })
    }
    return c.JSON(logs)
}