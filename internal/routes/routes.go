package routes

import (
	"Skripsigma-BE/internal/config"
	"Skripsigma-BE/internal/handler"
	"Skripsigma-BE/internal/middleware"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

// Setup route handler
func Setup(app *fiber.App) {

	// repositories
	userRepo := repository.NewUserRepository(config.DB)
	researchCaseRepo := repository.NewResearchCaseRepository(config.DB)
	companyRepo := repository.NewCompanyRepository(config.DB)
	tagRepo := repository.NewTagRepository(config.DB)
	roleRepo := repository.NewRoleRepository(config.DB)
	applicationRepo := repository.NewApplicationRepository(config.DB)

	// Services
	authService := service.NewAuthService(userRepo)
	researchCaseService := service.NewResearchCaseService(researchCaseRepo)
	companyService := service.NewCompanyService(companyRepo)
	tagService := service.NewTagService(tagRepo)
	roleService := service.NewRoleService(roleRepo)
	applicationService := service.NewApplicationService(applicationRepo, roleRepo, researchCaseRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	researchCaseHandler := handler.NewResearchCaseHandler(researchCaseService)
	companyHandler := handler.NewCompanyHandler(companyService)
	tagHandler := handler.NewTagHandler(tagService)
	roleHandler := handler.NewRoleHandler(roleService)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// Auth routes
	app.Post("/api/register", authHandler.Register)
	app.Post("/api/login", authHandler.Login)
	app.Get("/api/user", authMiddleware, authHandler.GetUserData)

	// Research Case routes
	app.Post("/api/research-case", researchCaseHandler.CreateResearchCase)
	app.Get("/api/research-case", researchCaseHandler.GetAllResearchCases)
	app.Get("/api/research-case/:id", researchCaseHandler.GetResearchCaseByID)

	// Company routes
	app.Post("/api/company", companyHandler.CreateCompany)
	app.Get("/api/company", companyHandler.GetAllCompanies)
	app.Get("/api/company/:id", companyHandler.GetCompanyByID)

	// Tag routes
	app.Post("/api/tag", tagHandler.CreateTag)
	app.Get("/api/tag", tagHandler.GetAllTags)
	app.Get("/api/tag/:id", tagHandler.GetTagByID)

	// Role routes
	app.Post("/api/role", roleHandler.CreateRole)
	app.Get("/api/role", roleHandler.GetAllRoles)
	app.Get("/api/role/:id", roleHandler.GetRoleByID)

	// Application routes
	app.Post("/api/application", authMiddleware, applicationHandler.CreateApplication)
	// app.Post("/api/application", middleware.StudentOnly(), applicationHandler.CreateApplication)
	// app.Get("/api/application", applicationHandler.GetAllApplications)
	// app.Get("/api/application/:id", applicationHandler.GetApplicationByID)
}
