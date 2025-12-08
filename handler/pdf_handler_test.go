package handler_test

import (
	"auth-service/handler"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockPdfService struct {
	ReturnError bool
}

func (m *MockPdfService) GenerateTripReceiptPDF(tripID string) ([]byte, string, error) {
	if tripID == "error" {
		return nil, "", errors.New("failed to generate PDF")
	}
	return []byte("PDF content"), "trip_receipt.pdf", nil
}

func TestHandlerPdfReceipt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		tripID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			tripID:         "123",
			expectedStatus: http.StatusOK,
			expectedBody:   "PDF content",
		},
		{
			name:           "service error",
			tripID:         "error",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to generate PDF"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			mockService := &MockPdfService{}
			h := handler.NewPDFHandler(mockService)
			r.GET("/pdf/:trip_id", h.HandlePDFReceipt)
			req, _ := http.NewRequest("GET", "/pdf/"+tt.tripID, nil)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.String())
		})
	}
}
