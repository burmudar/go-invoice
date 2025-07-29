package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/burmudar/go-invoice/db"
	"github.com/burmudar/go-invoice/eventbus"
)

type Service struct {
	http.Server
	store   *db.Store
	bus     *eventbus.Bus
	running atomic.Bool
}

func NewService(address string, store *db.Store, bus *eventbus.Bus) *Service {
	mux := http.NewServeMux()
	sv := &Service{
		store: store,
		bus:   bus,
	}

	mux.HandleFunc("/invoice", sv.handleNewInvoice)
	mux.HandleFunc("/invoice/{id}/details", sv.handleUpdateInvoice)
	mux.HandleFunc("/invoice/{id}/finalize", sv.handleFinalizeInvoice)
	mux.HandleFunc("/invoice/{id}/pdf", sv.handleInvoicePdf)
	mux.HandleFunc("/invoice/{id}/", sv.handleGetInvoice)

	sv.Handler = mux
	sv.Addr = address

	return sv
}

func (s *Service) handleNewInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "New invoice")

}
func (s *Service) handleUpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Update invoice [%s]", id)
}
func (s *Service) handleFinalizeInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Finalize invoice [%s]", id)
}
func (s *Service) handleInvoicePdf(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "PDF invoice [%s]", id)
}

func (s *Service) handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Get invoice [%s]", id)
}

func (s *Service) Start() error {
	if s.running.Load() {
		return nil
	}

	cnf := net.ListenConfig{}
	l, err := cnf.Listen(context.Background(), "tcp4", s.Addr)
	if err != nil {
		return err
	}

	go func() {
		s.running.Store(true)
		fmt.Printf("serving on address: %s\n", l.Addr())
		if err := s.Serve(l); err != nil {
			fmt.Printf("service listen error: %v\n", err)
		}
	}()

	return nil
}
