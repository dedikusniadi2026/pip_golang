package service

type PDFServiceInterface interface {
	GenerateTripReceiptPDF(tripID string) ([]byte, string, error)
}
