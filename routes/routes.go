package routes

import (
	"Skripsigma-BE/controllers"
	// companyController "Skripsigma-BE/controllers/companyController"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Auth routes
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	// Company routes
	company := app.Group("/api/companies")
	company.Post("/", controllers.CreateCompany)
	company.Get("/", controllers.GetAllCompanies)
	company.Get("/:id", controllers.GetCompanyByID)
	company.Put("/:id", controllers.UpdateCompany)
	company.Delete("/:id", controllers.DeleteCompany)

	// Research case routes
	researchCases := app.Group("/api/research-cases")
	researchCases.Post("/", controllers.CreateResearchCase)
	researchCases.Get("/", controllers.GetResearchCases)
	researchCases.Get("/:id", controllers.GetResearchCase)
	researchCases.Put("/:id", controllers.UpdateResearchCase)
	researchCases.Delete("/:id", controllers.DeleteResearchCase)

	// Tags routes
	tags := app.Group("/api/tags")
	tags.Post("/", controllers.CreateTag)
	tags.Get("/", controllers.GetTags)
	tags.Get("/:id", controllers.GetTag)
	tags.Put("/:id", controllers.UpdateTag)
	tags.Delete("/:id", controllers.DeleteTag)

	// // Research case tags routes
	// tags := app.Group("/api/research-case-tags")
	// tags.Post("/", researchCaseTagController.CreateResearchCaseTag)
	// tags.Get("/", researchCaseTagController.GetAllResearchCaseTags)
	// tags.Get("/:id", researchCaseTagController.GetResearchCaseTagByID)
	// tags.Put("/:id", researchCaseTagController.UpdateResearchCaseTag)
	// tags.Delete("/:id", researchCaseTagController.DeleteResearchCaseTag)
}
