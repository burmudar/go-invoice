package invoice

import (
	"context"
	"log/slog"

	"github.com/burmudar/go-invoice/internal/app/eventbus"
	"github.com/burmudar/go-invoice/internal/domain"
	"github.com/burmudar/go-invoice/internal/ports/inbound"
	"github.com/burmudar/go-invoice/internal/ports/outbound"
)

type Service struct {
	repo outbound.InvoiceRepository
	bus  *eventbus.Bus
}

func New(logger slog.Logger, repo outbound.InvoiceRepository, bus *eventbus.Bus) *Service {
	return &Service{
		repo: repo, bus: bus,
	}
}

// CreateNewInvoice implements inbound.InvoiceService.
func (s *Service) CreateNewInvoice(ctx context.Context, req inbound.CreateInvoiceRequest) (*domain.Invoice, error) {
	panic("unimplemented")
}

// FinalizeInvoice implements inbound.InvoiceService.
func (s *Service) FinalizeInvoice(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GeneratePDF implements inbound.InvoiceService.
func (s *Service) GeneratePDF(ctx context.Context, id string) ([]byte, error) {
	panic("unimplemented")
}

// GetInvoice implements inbound.InvoiceService.
func (s *Service) GetInvoice(ctx context.Context, id string) *domain.Invoice {
	panic("unimplemented")
}

// UpdateInvoice implements inbound.InvoiceService.
func (s *Service) UpdateInvoice(ctx context.Context, id string, req inbound.UpdateInvoiceRequest) error {
	panic("unimplemented")
}

var _ inbound.InvoiceService = (&Service{})
