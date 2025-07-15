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
	menuRepo := repository.NewMenuRepository(config.DB)
	chatRepo := repository.NewChatRepository(config.DB)
	notificationRepo := repository.NewNotificationRepository(config.DB)
	weeklyReportRepo := repository.NewWeeklyReportRepository(config.DB)
	universityRepo := repository.NewUniversityRepository(config.DB)
	userCreateRepo := repository.NewUserCreateRepository(config.DB)
	userCreateLogRepo := repository.NewUserCreateLogRepository(config.DB)
	assignmentRepo := repository.NewAssignmentRepository(config.DB)

	// Services
	authService := service.NewAuthService(userRepo)
	UserService := service.NewUserService(userRepo)
	researchCaseService := service.NewResearchCaseService(researchCaseRepo)
	companyService := service.NewCompanyService(companyRepo)
	tagService := service.NewTagService(tagRepo)
	roleService := service.NewRoleService(roleRepo)
	applicationService := service.NewApplicationService(applicationRepo, roleRepo, researchCaseRepo, userRepo, assignmentRepo)
	menuService := service.NewMenuService(menuRepo, assignmentRepo, config.DB)
	chatService := service.NewChatService(chatRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	weeklyReportService := service.NewWeeklyReportService(weeklyReportRepo)
	universityService := service.NewUniversityService(universityRepo)
	userCreateService := service.NewUserCreateService(userCreateRepo)
	userCreateLogService := service.NewUserCreateLogService(userCreateLogRepo)
	assignmentService := service.NewAssignmentService(assignmentRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(UserService)
	researchCaseHandler := handler.NewResearchCaseHandler(researchCaseService)
	companyHandler := handler.NewCompanyHandler(companyService)
	tagHandler := handler.NewTagHandler(tagService)
	roleHandler := handler.NewRoleHandler(roleService)
	applicationHandler := handler.NewApplicationHandler(applicationService, notificationService, researchCaseService)
	menuHandler := handler.NewMenuHandler(menuService)
	chatHandler := handler.NewChatHandler(chatService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	weeklyReportHandler := handler.NewWeeklyReportHandler(weeklyReportService)
	universityHandler := handler.NewUniversityHandler(universityService)
	userCreateHandler := handler.NewUserCreateHandler(userCreateService, userCreateLogService)
	assignmentHandler := handler.NewAssignmentHandler(assignmentService)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// Auth routes
	app.Post("/api/register", authHandler.Register)
	app.Post("/api/login", authHandler.Login)
	app.Get("/api/user", authMiddleware, authHandler.GetUserData)
	// logout
	app.Post("/api/logout", authMiddleware, authHandler.Logout)

	// User routes
	app.Put("/api/profile", authMiddleware, userHandler.UpdateProfile)
	app.Put("/api/profile/photo", authMiddleware, userHandler.UpdateProfilePhoto)

	// Research Case routes
	app.Post("/api/research-case", researchCaseHandler.CreateResearchCase)
	app.Get("/api/research-case", researchCaseHandler.GetAllResearchCases)
	app.Get("/api/research-case/:id", researchCaseHandler.GetResearchCaseByID)
	app.Get("/api/research-case/company/:company_id", researchCaseHandler.GetResearchCasesByCompanyID)
	app.Put("/api/research-case/:id", researchCaseHandler.UpdateResearchCase)

	// Company routes
	app.Post("/api/company", companyHandler.CreateCompany)
	app.Get("/api/company", companyHandler.GetAllCompanies)
	app.Get("/api/company/:id", companyHandler.GetCompanyByID)
	app.Put("/api/company/:id", companyHandler.UpdateCompany)

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
	app.Put("/api/application/review/:id", authMiddleware, applicationHandler.ProcessApplication)
	app.Put("/api/application/respond-to/:id", authMiddleware, applicationHandler.RespondToApplication)
	app.Get("/api/application/research-case/:id", authMiddleware, applicationHandler.GetApplicationsByResearchCaseID)
	app.Get("/api/application/student/:id", authMiddleware, applicationHandler.GetApplicationsByStudentID)
	// cek application exist
	app.Get("/api/application/check", authMiddleware, applicationHandler.CheckApplicationExists)
	// app.Post("/api/application", middleware.StudentOnly(), applicationHandler.CreateApplication)
	// app.Get("/api/application", applicationHandler.GetAllApplications)
	// app.Get("/api/application/:id", applicationHandler.GetApplicationByID)

	// menu routes
	app.Post("/api/menu", authMiddleware, menuHandler.CreateMenu)
	app.Get("/api/menu", authMiddleware, menuHandler.GetAllMenu)

	app.Get("/api/chatrooms/student/:id", authMiddleware, chatHandler.GetChatRoomsByStudentID)
	app.Get("/api/chatrooms/company/:id", authMiddleware, chatHandler.GetChatRoomByCompanyID)
	app.Get("/api/chatrooms/messages/:room_id", authMiddleware, chatHandler.GetMessagesByRoomID)

	// Notification routes
	app.Get("/api/notifications", authMiddleware, notificationHandler.GetNotifications)
	app.Post("/api/notifications", authMiddleware, notificationHandler.CreateNotification)
	app.Put("/api/notifications/:id/true", authMiddleware, notificationHandler.MarkAsRead)
	app.Put("/api/notifications/read-all", authMiddleware, notificationHandler.MarkAllAsRead)

	// Weekly Report routes
	app.Post("/api/weekly-report", authMiddleware, weeklyReportHandler.SubmitWeeklyReport)
	app.Get("/api/weekly-report", authMiddleware, weeklyReportHandler.GetWeeklyReports)
	app.Get("/api/weekly-report/bysupervisor", authMiddleware, weeklyReportHandler.GetWeeklyReportsForSupervisor)

	// university
	app.Get("/api/university", authMiddleware, universityHandler.GetAll)

	// create user
	app.Post("/api/user-create", authMiddleware, userCreateHandler.CreateUser)
	app.Get("/api/user-create", authMiddleware, userCreateHandler.GetCreationLogs)

	// assignment
	app.Get("/api/assignment/active", authMiddleware, assignmentHandler.GetMyActiveAssignment)


}
