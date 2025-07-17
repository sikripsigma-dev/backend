package service

import (
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/repository"
	"Skripsigma-BE/internal/util"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type StudentDocumentService interface {
	Create(userID, docType string, file *multipart.FileHeader) (models.StudentDocument, error)
	GetByUserID(userID string) ([]models.StudentDocument, error)
}

type studentDocumentService struct {
	studentDocumentRepo repository.StudentDocumentRepository
	userRepo            repository.UserRepository
}

func NewStudentDocumentService(
	studentDocumentRepo repository.StudentDocumentRepository,
	userRepo repository.UserRepository,
) StudentDocumentService {
	return &studentDocumentService{
		studentDocumentRepo: studentDocumentRepo,
		userRepo:            userRepo,
	}
}

func (s *studentDocumentService) Create(userID, docType string, file *multipart.FileHeader) (models.StudentDocument, error) {
	var doc models.StudentDocument

	// Validasi user
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return doc, fmt.Errorf("user tidak ditemukan: %w", err)
	}

	// Simpan file ke folder
	dir := "./public/documents/student"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return doc, fmt.Errorf("gagal membuat direktori: %w", err)
		}
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%s%s", userID, docType, ext)
	filePath := filepath.Join(dir, fileName)

	// Hapus file lama dengan pattern userID_type.*
	pattern := filepath.Join(dir, fmt.Sprintf("%s_%s.*", userID, docType))
	matches, _ := filepath.Glob(pattern)
	for _, match := range matches {
		_ = os.Remove(match)
	}

	if err := util.SaveUploadedFile(file, filePath); err != nil {
		return doc, fmt.Errorf("gagal menyimpan file: %w", err)
	}

	doc = models.StudentDocument{
		ID: uuid.New().String(),
		UserID: userID,
		Type:   docType,
		Path:   "/documents/student/" + fileName,
	}

	if err := s.studentDocumentRepo.Save(&doc); err != nil {
		return doc, fmt.Errorf("gagal menyimpan ke database: %w", err)
	}

	return doc, nil
}

func (s *studentDocumentService) GetByUserID(userID string) ([]models.StudentDocument, error) {
	docs, err := s.studentDocumentRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("document(s) not found")
	}
	return docs, nil
}
