package main

import (
	"time"

	"github.com/burmudar/go-invoice/internal/adapters/httpserver"
)

func main() {
	svc = app.
	handlers = httpserver.InvoiceHandler{

	}
	httpserver.NewServer()
	sv := api.NewService("127.0.0.1:11111", nil, nil)

	sv.Start()

	<-time.After(30 * time.Second)
}
