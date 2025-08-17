package outbound

import (
	"context"

	"github.com/burmudar/go-invoice/internal/domain"
)

type InvoiceRepository interface {
	Create(ctx context.Context, i *domain.Invoice) (*domain.Invoice, error)
	Update(ctx context.Context, i *domain.Invoice) error
	Get(ctx context.Context, id domain.InvoiceID) (*domain.Invoice, error)
	Delete(ctx context.Context, id domain.InvoiceID) error
}
