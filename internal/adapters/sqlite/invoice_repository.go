package sqlite

import (
	"context"
	"time"

	"github.com/burmudar/go-invoice/internal/domain"
	"github.com/burmudar/go-invoice/internal/ports/outbound"
	"gorm.io/gorm"
)

type InvoiceRecord struct {
	// Model includes ID, CreatedAt, UpdatedAt and DeletedAt
	gorm.Model
	ID          string
	CustomerID  string
	LineItems   []LineItem
	Status      string
	CreatedAt   time.Time
	FinalizedAt time.Time
	TotalAmount float64
	Currency    string
}

type LineItem struct {
	gorm.Model
	ID          string
	InvoiceID   string
	Description string
	Qty         int
	Amount      float64
}

type InvoiceRepository struct {
	db *gorm.DB
}

func toInvoiceRecord(i *domain.Invoice) *InvoiceRecord {
	return &InvoiceRecord{}
}

func toInvoiceDomain(i *InvoiceRecord) *domain.Invoice {
	return &domain.Invoice{}
}

func (repo *InvoiceRepository) handle() gorm.Interface[InvoiceRecord] {
	return gorm.G[InvoiceRecord](repo.db)
}

// Create implements outbound.InvoiceRepository.
func (repo *InvoiceRepository) Create(ctx context.Context, i *domain.Invoice) (*domain.Invoice, error) {
	rec := toInvoiceRecord(i)
	err := repo.handle().Create(ctx, rec)
	return toInvoiceDomain(rec), err
}

// Delete implements outbound.InvoiceRepository.
func (i *InvoiceRepository) Delete(ctx context.Context, id domain.InvoiceID) error {
	panic("unimplemented")
}

// Get implements outbound.InvoiceRepository.
func (i *InvoiceRepository) Get(ctx context.Context, id domain.InvoiceID) (*domain.Invoice, error) {
	panic("unimplemented")
}

// Update implements outbound.InvoiceRepository.
func (*InvoiceRepository) Update(ctx context.Context, i *domain.Invoice) error {
	panic("unimplemented")
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{
		db: db,
	}
}

var _ outbound.InvoiceRepository = &InvoiceRepository{}
