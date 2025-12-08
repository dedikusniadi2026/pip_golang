package handler

import (
	"auth-service/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PDFHandler struct {
	PDFService service.PDFServiceInterface
}

func NewPDFHandler(s service.PDFServiceInterface) *PDFHandler {
	return &PDFHandler{PDFService: s}
}

func (h *PDFHandler) HandlePDFReceipt(c *gin.Context) {
	tripID := c.Param("trip_id")

	pdfBytes, filename, err := h.PDFService.GenerateTripReceiptPDF(tripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
