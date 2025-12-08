package service

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
)

type PaymentServiceInterface interface {
	GetPayments(page, pageSize int) ([]model.Payment, error)
	GetPaymentByID(ctx context.Context, id int) (*model.Payment, error)
	CreatePayment(ctx context.Context, p *model.Payment) (int, error)
	UpdatePayment(ctx context.Context, p *model.Payment) error
	DeletePayment(ctx context.Context, id int) error
	GetPaymentStats(ctx context.Context) (*model.PaymentStats, error)
}

type PaymentService struct {
	Repo repository.PaymentRepositoryInterface
}

func NewPaymentService(repo repository.PaymentRepositoryInterface) *PaymentService {
	return &PaymentService{
		Repo: repo,
	}
}

func (s *PaymentService) GetPayments(page, pageSize int) ([]model.Payment, error) {
	return s.Repo.GetPayments(page, pageSize)
}

func (s *PaymentService) GetPaymentStats(ctx context.Context) (*model.PaymentStats, error) {
	return s.Repo.GetPaymentStats(ctx)
}

func (s *PaymentService) GetAll(ctx context.Context) ([]model.Payment, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (*model.Payment, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment *model.Payment) (int, error) {
	return s.Repo.Create(ctx, payment)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, p *model.Payment) error {
	return s.Repo.Update(ctx, p)
}

func (s *PaymentService) DeletePayment(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
