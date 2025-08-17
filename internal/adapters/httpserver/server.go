package httpserver

import "net/http"

func NewServer(address string, h *InvoiceHandler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/invoice", h.HandleNewInvoice)
	mux.HandleFunc("/invoice/{id}/details", h.HandleUpdateInvoice)
	mux.HandleFunc("/invoice/{id}/finalize", h.HandleFinalizeInvoice)
	mux.HandleFunc("/invoice/{id}/pdf", h.HandleGeneratePDF)
	mux.HandleFunc("/invoice/{id}/", h.HandleGetInvoice)

	sv := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return sv
}
