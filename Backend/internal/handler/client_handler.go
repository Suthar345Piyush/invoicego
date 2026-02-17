package handler

import "github.com/Suthar345Piyush/invoicego/internal/service"

// handler

type ClientHandler struct {
	clientService *service.ClientService
}

// client handler

func NewClientHandler(clientService *service.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

// creating client function
