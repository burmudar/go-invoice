package domain

import "time"

type InvoiceStatus string
type InvoiceID string

const DefaultCurrency = "USD"

type Money struct {
	Amount   float64
	Currency string
}

type Invoice struct {
	ID          InvoiceID
	CustomerID  string
	LineItems   []LineItem
	Status      InvoiceStatus
	CreatedAt   time.Time
	FinalizedAt time.Time
	Total       Money
}
type LineItem struct {
	ID          string
	InvoiceID   InvoiceID
	Description string
	Qty         int
	Amount      float64
}

type Customer struct {
	ID      string
	Name    string
	Address string
}

func NewInvoice() *Invoice {
	return &Invoice{
		LineItems: make([]LineItem, 0),
	}
}

func (i *Invoice) AddLineItem(item LineItem) error {
	item.InvoiceID = i.ID
	i.LineItems = append(i.LineItems, item)

	return nil
}

func (i *Invoice) CalculateTotal() Money {
	total := Money{
		Amount:   0.00,
		Currency: DefaultCurrency,
	}
	for _, item := range i.LineItems {
		total.Amount += item.Total()
	}

	return total
}

func (l *LineItem) Total() float64 {
	return l.Amount * float64(l.Qty)
}
