package main

import (
	"time"

	"github.com/burmudar/go-invoice/api"
)

func main() {
	sv := api.NewService("127.0.0.1:11111", nil, nil)

	sv.Start()

	<-time.After(30 * time.Second)
}
