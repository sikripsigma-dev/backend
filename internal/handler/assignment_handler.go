package handler

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

type AssignmentHandler struct {
	assignmentService *service.AssignmentService
}

func NewAssignmentHandler(assignmentService *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{assignmentService: assignmentService}
}

// GET /api/assignment/active (by student)
func (h *AssignmentHandler) GetMyActiveAssignment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	assignment, err := h.assignmentService.GetActiveAssignment(user.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No active assignment found"})
	}
	return c.JSON(fiber.Map{"data": assignment})
}

// GET /api/assignment/company (by company)
func (h *AssignmentHandler) GetActiveAssignmentsByCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	companyID := user.Company.CompanyID

	assignments, err := h.assignmentService.GetActiveAssignmentsByCompany(companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data assignment"})
	}
	return c.JSON(fiber.Map{"data": assignments})
}

// PUT /api/assignment/:id/status
func (h *AssignmentHandler) UpdateAssignmentStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignment ID"})
	}
	if err := h.assignmentService.UpdateAssignmentStatus(uint(id), req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Assignment status updated"})
}

// GET /api/assignments (list assignments milik user login)
func (h *AssignmentHandler) ListMyAssignments(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	res, err := h.assignmentService.GetAssignmentsByUser(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": res})
}

// GET /api/assignments/:id/certificate (download PDF certificate)
func (h *AssignmentHandler) GenerateCertificate(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	idUint, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id"})
	}

	// pastikan assignment milik user
	a, err := h.assignmentService.GetOwnedAssignment(user.Id, uint(idUint))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Assignment tidak ditemukan"})
	}

	if a.Status != "completed" || a.EndedAt == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Sertifikat hanya tersedia untuk assignment yang sudah selesai",
		})
	}

	// === Generate PDF ===
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Border
	pdf.SetDrawColor(218, 165, 32)
	pdf.SetLineWidth(2)
	pdf.Rect(5, 5, 287, 200, "D")

	// Background
	pdf.SetFillColor(255, 250, 240)
	pdf.Rect(7, 7, 283, 196, "F")

	// Title
	pdf.SetFont("Times", "B", 28)
	pdf.SetTextColor(0, 51, 102)
	pdf.CellFormat(0, 30, "SERTIFIKAT PENYELESAIAN", "", 1, "C", false, 0, "")

	// Sub title
	pdf.SetFont("Times", "I", 16)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 10, "Diberikan kepada:", "", 1, "C", false, 0, "")

	// Nama
	pdf.SetFont("Times", "B", 26)
	pdf.SetTextColor(139, 0, 0)
	pdf.Ln(3)
	pdf.CellFormat(0, 15, user.Name, "", 1, "C", false, 0, "")

	// Deskripsi
	pdf.SetFont("Times", "", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "Telah menyelesaikan studi kasus:", "", 1, "C", false, 0, "")

	pdf.SetFont("Times", "B", 18)
	pdf.SetTextColor(0, 0, 128)
	pdf.Ln(3)
	title := a.ResearchCase.Title
	if title == "" {
		title = "—"
	}
	pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")

	// Periode
	pdf.SetFont("Times", "", 14)
	pdf.SetTextColor(50, 50, 50)
	pdf.Ln(8)
	endStr := "-"
	if a.EndedAt != nil {
		endStr = a.EndedAt.Format("02 January 2006")
	}
	periode := "Periode: " + a.StartedAt.Format("02 January 2006") + " - " + endStr
	pdf.CellFormat(0, 10, periode, "", 1, "C", false, 0, "")

	// TTD
	pdf.Ln(25)
	pdf.SetFont("Times", "I", 14)
	pdf.CellFormat(100, 10, "Dosen Pembimbing", "", 0, "C", false, 0, "")
	pdf.CellFormat(87, 10, "", "", 0, "C", false, 0, "")
	pdf.CellFormat(100, 10, "Perwakilan Perusahaan", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate certificate",
			"error":   err.Error(),
		})
	}

	c.Set("Content-Type", "application/pdf")
	filename := fmt.Sprintf("Sertifikat_%s_%s.pdf", title, time.Now().Format("2006-01-02"))
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.Send(buf.Bytes())
}

// func (h *AssignmentHandler) GenerateCertificate(c *fiber.Ctx) error {
//     user := c.Locals("user").(*models.User)
//     userID := user.Id

//     assignment, err := h.assignmentService.GetActiveAssignment(userID)
//     if err != nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No active assignment found"})
//     }

//     if assignment.Status != "completed" || assignment.EndedAt == nil {
//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//             "message": "Sertifikat hanya tersedia untuk assignment yang sudah selesai",
//             "details": fiber.Map{
//                 "required_status": "completed",
//                 "current_status": assignment.Status,
//                 "has_end_date": assignment.EndedAt != nil,
//             },
//         })
//     }

//     // === Generate PDF ===
//     pdf := gofpdf.New("L", "mm", "A4", "")
//     pdf.AddPage()

//     // --- Border luar (emas) ---
//     pdf.SetDrawColor(218, 165, 32) // gold
//     pdf.SetLineWidth(2)
//     pdf.Rect(5, 5, 287, 200, "D")

//     // --- Background cream ---
//     pdf.SetFillColor(255, 250, 240) // Floral White
//     pdf.Rect(7, 7, 283, 196, "F")

//     // --- Logo (opsional, pastikan path valid) ---
//     // pdf.Image("./assets/logo.png", 20, 10, 40, 0, false, "", 0, "")

//     // --- Judul ---
//     pdf.SetFont("Times", "B", 28)
//     pdf.SetTextColor(0, 51, 102) // Navy Blue
//     pdf.CellFormat(0, 30, "SERTIFIKAT PENYELESAIAN", "", 1, "C", false, 0, "")

//     // --- Sub Judul ---
//     pdf.SetFont("Times", "I", 16)
//     pdf.SetTextColor(80, 80, 80)
//     pdf.CellFormat(0, 10, "Diberikan kepada:", "", 1, "C", false, 0, "")

//     // --- Nama Peserta ---
//     pdf.SetFont("Times", "B", 26)
//     pdf.SetTextColor(139, 0, 0) // Dark Red
//     pdf.Ln(3)
//     pdf.CellFormat(0, 15, user.Name, "", 1, "C", false, 0, "")

//     // --- Deskripsi ---
//     pdf.SetFont("Times", "", 16)
//     pdf.SetTextColor(0, 0, 0)
//     pdf.Ln(5)
//     pdf.CellFormat(0, 10, "Telah menyelesaikan studi kasus:", "", 1, "C", false, 0, "")

//     pdf.SetFont("Times", "B", 18)
//     pdf.SetTextColor(0, 0, 128)
//     pdf.Ln(3)
//     pdf.CellFormat(0, 10, assignment.ResearchCase.Title, "", 1, "C", false, 0, "")

//     // --- Periode ---
//     pdf.SetFont("Times", "", 14)
//     pdf.SetTextColor(50, 50, 50)
//     pdf.Ln(8)
//     periode := "Periode: " + assignment.StartedAt.Format("02 January 2006") +
//         " - " + assignment.EndedAt.Format("02 January 2006")
//     pdf.CellFormat(0, 10, periode, "", 1, "C", false, 0, "")

//     // --- Space tanda tangan ---
//     pdf.Ln(25)
//     pdf.SetFont("Times", "I", 14)
//     pdf.CellFormat(100, 10, "Dosen Pembimbing", "", 0, "C", false, 0, "")
//     pdf.CellFormat(87, 10, "", "", 0, "C", false, 0, "")
//     pdf.CellFormat(100, 10, "Perwakilan Perusahaan", "", 1, "C", false, 0, "")

//     // Bisa tambahkan tanda tangan digital
//     // pdf.Image("./assets/signature.png", 50, 150, 40, 0, false, "", 0, "")
//     // pdf.Image("./assets/stamp.png", 200, 145, 35, 0, false, "", 0, "")

//     // === Output PDF ===
//     var buf bytes.Buffer
//     err = pdf.Output(&buf)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "message": "Failed to generate certificate",
//             "error":   err.Error(),
//         })
//     }

//     c.Set("Content-Type", "application/pdf")
//     filename := fmt.Sprintf("Sertifikat_%s_%s.pdf", assignment.ResearchCase.Title, time.Now().Format("2006-01-02"))
//     c.Set("Content-Disposition", "attachment; filename="+filename)

//     return c.Send(buf.Bytes())
// }
