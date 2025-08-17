package inbound

import (
	"context"

	"github.com/burmudar/go-invoice/internal/domain"
)

type CreateInvoiceRequest struct {
	CustomerID string
	LineItems  []domain.LineItem
}
type UpdateInvoiceRequest struct {
	ID        domain.InvoiceID
	LineItems []domain.LineItem
}

type InvoiceService interface {
	CreateNewInvoice(ctx context.Context, req CreateInvoiceRequest) (*domain.Invoice, error)
	UpdateInvoice(ctx context.Context, id string, req UpdateInvoiceRequest) error
	FinalizeInvoice(ctx context.Context, id string) error
	GetInvoice(ctx context.Context, id string) *domain.Invoice
	GeneratePDF(ctx context.Context, id string) ([]byte, error)
}
