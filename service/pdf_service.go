package service

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
	"fmt"
	"net/url"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFGenerator interface {
	GeneratePDF(html string) ([]byte, error)
}

type DefaultPDFGenerator struct{}

func (d *DefaultPDFGenerator) GeneratePDF(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfBytes []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+url.PathEscape(html)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBytes, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	)

	return pdfBytes, err
}

type PDFServiceInterface interface {
	GenerateTripReceiptPDF(tripID string) ([]byte, string, error)
}

type TemplateRenderer interface {
	RenderPDFReceipt(data model.PDFTemplateData) (string, error)
}

type PDFService struct {
	TripRepo         repository.PDFRepositoryInterface
	PDFGenerator     PDFGenerator
	TemplateRenderer TemplateRenderer
}

func (s *PDFService) GenerateTripReceiptPDF(tripID string) ([]byte, string, error) {
	tripData, err := s.TripRepo.GetTripByID(tripID)
	if err != nil {
		return nil, "", fmt.Errorf("gagal mendapatkan data trip: %w", err)
	}

	amountFormatted := fmt.Sprintf("%.0f", float64(tripData.Amount))

	html, err := s.TemplateRenderer.RenderPDFReceipt(model.PDFTemplateData{
		Pdf:             tripData,
		AmountFormatted: amountFormatted,
	})
	if err != nil {
		return nil, "", fmt.Errorf("gagal render template: %w", err)
	}

	pdfBytes, err := s.PDFGenerator.GeneratePDF(html)
	if err != nil {
		return nil, "", fmt.Errorf("gagal generate PDF: %w", err)
	}

	filename := fmt.Sprintf("receipt_%s.pdf", tripID)
	return pdfBytes, filename, nil
}
