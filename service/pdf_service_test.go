package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPDFRepository struct{ mock.Mock }
type MockPDFGenerator struct{ mock.Mock }
type MockTemplateRenderer struct{ mock.Mock }

func (m *MockPDFRepository) GetTripByID(tripID string) (model.Pdf, error) {
	args := m.Called(tripID)
	return args.Get(0).(model.Pdf), args.Error(1)
}

func (m *MockPDFGenerator) GeneratePDF(html string) ([]byte, error) {
	args := m.Called(html)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockTemplateRenderer) RenderPDFReceipt(data model.PDFTemplateData) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func TestPDFService_GenerateTripReceiptPDF(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	tripPdf := model.Pdf{
		ID:           "123",
		Amount:       100000,
		CustomerName: "John Doe",
	}

	mockRepo.On("GetTripByID", "trip123").Return(tripPdf, nil)

	mockTpl.On("RenderPDFReceipt", mock.AnythingOfType("model.PDFTemplateData")).
		Return("<html>dummy</html>", nil)

	pdfBytesDummy := []byte("dummy pdf content")
	mockGen.On("GeneratePDF", "<html>dummy</html>").Return(pdfBytesDummy, nil)

	pdfBytes, filename, err := svc.GenerateTripReceiptPDF("trip123")

	assert.NoError(t, err)
	assert.Equal(t, pdfBytesDummy, pdfBytes)
	assert.Equal(t, "receipt_trip123.pdf", filename)

	mockRepo.AssertExpectations(t)
	mockGen.AssertExpectations(t)
	mockTpl.AssertExpectations(t)
}

func TestPDFService_GenerateTripReceiptPDF_GetTrip_Error(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	mockRepo.On("GetTripByID", "trip123").
		Return(model.Pdf{}, assert.AnError)

	pdfBytes, filename, err := svc.GenerateTripReceiptPDF("trip123")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "", filename)

	mockRepo.AssertExpectations(t)
}

func TestPDFService_GenerateTripReceiptPDF_Template_Error(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	tripPdf := model.Pdf{ID: "123"}
	mockRepo.On("GetTripByID", "trip123").Return(tripPdf, nil)

	mockTpl.On("RenderPDFReceipt", mock.AnythingOfType("model.PDFTemplateData")).
		Return("", assert.AnError)

	pdfBytes, filename, err := svc.GenerateTripReceiptPDF("trip123")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "", filename)

	mockRepo.AssertExpectations(t)
	mockTpl.AssertExpectations(t)
}

func TestPDFService_GenerateTripReceiptPDF_GeneratePDF_Error(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	tripPdf := model.Pdf{ID: "123"}
	mockRepo.On("GetTripByID", "trip123").Return(tripPdf, nil)

	mockTpl.On("RenderPDFReceipt", mock.AnythingOfType("model.PDFTemplateData")).
		Return("<html>dummy</html>", nil)

	mockGen.On("GeneratePDF", "<html>dummy</html>").Return([]byte(nil), assert.AnError)

	pdfBytes, filename, err := svc.GenerateTripReceiptPDF("trip123")

	assert.Error(t, err)
	assert.Nil(t, pdfBytes)
	assert.Equal(t, "", filename)

	mockRepo.AssertExpectations(t)
	mockTpl.AssertExpectations(t)
	mockGen.AssertExpectations(t)
}

func TestPDFService_RepoError(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	mockRepo.On("GetTripByID", "trip123").Return(model.Pdf{}, errors.New("repo error"))

	_, _, err := svc.GenerateTripReceiptPDF("trip123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gagal mendapatkan data trip")
}

func TestPDFService_TemplateError(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	tripPdf := model.Pdf{ID: "123", Amount: 100000, CustomerName: "John"}
	mockRepo.On("GetTripByID", "trip123").Return(tripPdf, nil)
	mockTpl.On("RenderPDFReceipt", mock.Anything).Return("", errors.New("template error"))

	_, _, err := svc.GenerateTripReceiptPDF("trip123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gagal render template")
}

func TestPDFService_PDFGeneratorError(t *testing.T) {
	mockRepo := new(MockPDFRepository)
	mockGen := new(MockPDFGenerator)
	mockTpl := new(MockTemplateRenderer)

	svc := &service.PDFService{
		TripRepo:         mockRepo,
		PDFGenerator:     mockGen,
		TemplateRenderer: mockTpl,
	}

	tripPdf := model.Pdf{ID: "123", Amount: 100000, CustomerName: "John"}
	mockRepo.On("GetTripByID", "trip123").Return(tripPdf, nil)
	mockTpl.On("RenderPDFReceipt", mock.Anything).Return("<html>dummy</html>", nil)
	mockGen.On("GeneratePDF", "<html>dummy</html>").Return([]byte(nil), errors.New("pdf gen error"))

	_, _, err := svc.GenerateTripReceiptPDF("trip123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gagal generate PDF")
}
