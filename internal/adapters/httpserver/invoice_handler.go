package httpserver

import (
	"fmt"
	"net/http"

	"github.com/burmudar/go-invoice/internal/ports/inbound"
)

type InvoiceHandler struct {
	svc inbound.InvoiceService
}

func (s *InvoiceHandler) HandleNewInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "New invoice")

}
func (s *InvoiceHandler) HandleUpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Update invoice [%s]", id)
}
func (s *InvoiceHandler) HandleFinalizeInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Finalize invoice [%s]", id)
}
func (s *InvoiceHandler) HandleGeneratePDF(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "PDF invoice [%s]", id)
}

func (s *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Get invoice [%s]", id)
}
