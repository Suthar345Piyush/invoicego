package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/Suthar345Piyush/invoicego/internal/middleware"
	"github.com/Suthar345Piyush/invoicego/internal/service"
	"github.com/Suthar345Piyush/invoicego/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// handler

type ClientHandler struct {
	clientService *service.ClientService
}

// client handler

func NewClientHandler(clientService *service.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

// creating client function

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var req domain.CreateClientRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	// validating the input

	if err := util.ValidateStruct(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	client, err := h.clientService.CreateClient(claims.UserID, &req)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	util.WriteSuccess(w, http.StatusCreated, client, "Client created successfully")
}

// function to get the client

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	clientIDStr := chi.URLParam(r, "id")

	//  parsing the client id string

	clientID, err := uuid.Parse(clientIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid client ID"))
	}

	client, err := h.clientService.GetClientByID(claims.UserID, clientID)

	if err != nil {
		util.WriteError(w, http.StatusNotFound, err)
	}

	util.WriteSuccess(w, http.StatusOK, client, "Client retrieved successfully")
}

// listing all the clients

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// parsing  pagination parameters

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	clients, err := h.clientService.GetClientsByUserID(claims.UserID, page, pageSize)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, clients, "Clients retrieved successfully")

}

// updating the client

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
	}

	clientIDStr := chi.URLParam(r, "id")

	clientID, err := uuid.Parse(clientIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid client ID"))
		return
	}

	// requesting to update the client

	var req domain.UpdateClientRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	// validating input

	if err := util.ValidateStruct(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// to update we need userID , clientID and req address

	client, err := h.clientService.UpdateClient(claims.UserID, clientID, &req)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}

	util.WriteSuccess(w, http.StatusOK, client, "Client updated successfully")
}

// delete client function

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	clientIDStr := chi.URLParam(r, "id")

	clientID, err := uuid.Parse(clientIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid client ID"))
		return
	}

	//deleting the client

	err = h.clientService.DeleteClient(claims.UserID, clientID)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, nil, "Client deleted successfully")

}
