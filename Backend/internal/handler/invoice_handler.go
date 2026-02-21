// invoice handler part

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

// invoice handler struct

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
	pdfService     *service.PDFService
	userService    *service.UserService
}

// invoice handler function

func NewInvoiceHandler(invoiceService *service.InvoiceService, pdfService *service.PDFService, userService *service.UserService) *InvoiceHandler {

	return &InvoiceHandler{
		invoiceService: invoiceService,
		pdfService:     pdfService,
		userService:    userService,
	}

}

//function for creating the invoice

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {

	// firstly getting the context of the user before creating invoice

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// getting invoice creation request

	var req domain.CreateInvoiceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	// validating the input

	if err := util.ValidateStruct(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	invoice, err := h.invoiceService.CreateInvoice(claims.UserID, &req)

	// status forbidden 403 - means , server understood the request , but refused to process it

	if err != nil {
		if err == domain.ErrInvoiceLimitExceeded {
			util.WriteError(w, http.StatusForbidden, err)
			return
		}

		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteSuccess(w, http.StatusCreated, invoice, "Invoice created successfully")

}

// function to get the invoices

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {

	// it will give the userID

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// invoice id string

	invoiceIDStr := chi.URLParam(r, "id")

	// parsing the invoice id string

	invoiceID, err := uuid.Parse(invoiceIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid invoice ID"))
		return
	}

	invoice, err := h.invoiceService.GetInvoiceByID(claims.UserID, invoiceID)

	if err != nil {
		util.WriteError(w, http.StatusNotFound, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, invoice, "Invoice retrieved successfully")

}

// list of invoices function

func (h *InvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// listing multiple invoices needs pagination and filter parameters

	// parsing pagination and filter parameters

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	status := r.URL.Query().Get("status")

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	invoices, err := h.invoiceService.GetInvoiceByUserID(claims.UserID, page, pageSize, status)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, invoices, "Invoice retrieved successfully")

}

// function for updating the invoice status

func (h *InvoiceHandler) UpdateInvoiceStatus(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	invoiceIDStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(invoiceIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid invoice ID"))
		return
	}

	// updating the invoice status

	var req domain.UpdateInvoiceStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	// validating the input

	if err := util.ValidateStruct(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// from invoice service , updating the status of the invoice

	invoice, err := h.invoiceService.UpdateInvoiceStatus(claims.UserID, invoiceID, &req)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, invoice, "Invoice status updated successfully")

}

// invoice delete function

func (h *InvoiceHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	invoiceIDStr := chi.URLParam(r, "id")
	// parsing the invoice id string
	invoiceID, err := uuid.Parse(invoiceIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid invoice ID"))
		return
	}

	err = h.invoiceService.DeleteInvoice(claims.UserID, invoiceID)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, nil, "Invoice deleted successfully")

}

// pdf generation function

func (h *InvoiceHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	invoiceIDStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(invoiceIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid invoice ID"))
		return
	}

	// getting invoices with full details

	invoice, err := h.invoiceService.GetInvoiceByID(claims.UserID, invoiceID)

	if err != nil {
		util.WriteError(w, http.StatusNotFound, err)
		return
	}

	// getting user details

	user, err := h.userService.GetUserByID(invoice.ClientID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// generate pdf method

	pdfBytes, err := h.pdfService.GenerateInvoicePDF(invoice, user)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// setting headers for pdf download

	w.Header().Set("Content-type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+invoice.InvoiceNumber+".pdf")
	w.Header().Set("Content-length", strconv.Itoa(len(pdfBytes)))

	// write pdf to response

	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)

}

// duplicate invoice function

func (h *InvoiceHandler) DuplicateInvoice(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	invoiceIDStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(invoiceIDStr)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, errors.New("invalid invoice ID"))
		return
	}

	invoice, err := h.invoiceService.DuplicateInvoice(claims.UserID, invoiceID)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteSuccess(w, http.StatusCreated, invoice, "Invoice duplicated successfully")

}

// invoice stats function
// using stats sql Query from invoice service into this get stats handler function

func (h *InvoiceHandler) GetStats(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		util.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	stats, err := h.invoiceService.GetInvoiceStats(claims.UserID)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, stats, "Invoice Stats retrieved successfully")

}
