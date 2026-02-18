package domain

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID                uuid.UUID      `json:"id"`
	UserID            uuid.UUID      `json:"user_id"`
	ClientID          uuid.UUID      `json:"client_id"`
	InvoiceNumber     string         `json:"invoice_number"`
	Status            string         `json:"status"`
	IssueDate         time.Time      `json:"issue_date"`
	DueDate           time.Time      `json:"due_date"`
	PaidDate          *time.Time     `json:"paid_date,omitempty"`
	Currency          string         `json:"currency"`
	Subtotal          float64        `json:"subtotal"`
	TaxRate           float64        `json:"tax_rate"`
	TaxAmount         float64        `json:"tax_amount"`
	DiscountAmount    float64        `json:"discount_amount"`
	TotalAmount       float64        `json:"total_amount"`
	TemplateID        string         `json:"template_id"`
	Notes             *string        `json:"notes,omitempty"`
	TermAndConditions *string        `json:"terms_and_conditions,omitempty"`
	PDFURL            *string        `json:"pdf_url,omitempty"`
	PDFGeneratedAt    *time.Time     `json:"pdf_generated_at,omitempty"`
	EmailSent         bool           `json:"email_sent"`
	EmailSentAt       *time.Time     `json:"email_sent_at,omitempty"`
	EmailOpened       bool           `json:"email_opened"`
	EmailOpenedAt     *time.Time     `json:"email_opened_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	Items             []*InvoiceItem `json:"items,omitempty"`
	Client            *Client        `json:"client,omitempty"`
}

type InvoiceItem struct {
	ID          uuid.UUID `json:"id"`
	InvoiceID   uuid.UUID `json:"invoice_id"`
	Description string    `json:"description"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Amount      float64   `json:"amount"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInvoiceRequest struct {
	ClientID          uuid.UUID               `json:"client_id" validate:"required"`
	IssueDate         string                  `json:"issue_date" validate:"required"`
	DueDate           string                  `json:"due_date" validate:"required"`
	Currency          string                  `json:"currency" validate:"required, len=3"`
	TaxRate           float64                 `json:"tax_rate" validate:"gte=0, lte=100"`
	DiscountAmount    float64                 `json:"discount_amount" validate:"gte-0"`
	TemplateID        string                  `json:"template_id"`
	Notes             *string                 `json:"notes,omitempty"`
	TermAndConditions *string                 `json:"terms_and_conditions,omitempty"`
	Items             []*CreateInvoiceItemReq `json:"items" validation:"required, min=1, dive"`
}

type CreateInvoiceItemReq struct {
	Description string  `json:"description" validate:"required"`
	Quantity    float64 `json:"quantity" validtate:"required, gte=0"`
	UnitPrice   float64 `json:"unit_price" validate:"required , gte=0"`
}

type UpdateInvoiceRequest struct {
	ClientID           *uuid.UUID              `json:"client_id,omitmepty"`
	IssueDate          *string                 `json:"issue_date,omitempty"`
	DueDate            *string                 `json:"due_date,omitempty"`
	Currency           *string                 `json:"currency,omitempty" validate:"omitempty,len=3"`
	TaxRate            *float64                `json:"tax_rate,omitempty" validate:"omitempty,gte=0,lte=100"`
	DiscountAmount     *float64                `json:"discount_amount,omitempty" validate:"omitempty,gte=0"`
	TemplateID         *string                 `json:"template_id,omitempty"`
	Notes              *string                 `json:"notes,omitempty"`
	TermsAndConditions *string                 `json:"terms_and_conditions,omitempty"`
	Items              []*CreateInvoiceItemReq `json:"items,omitempty" validate:"omitempty,min=1,dive"`
}

type UpdateInvoiceStatusRequest struct {
	Status   string  `json:"status" validate:"requrired, oneof=draft sent paid overdue canceled"`
	PaidDate *string `json:"paid_date,omitempty"`
}

type InvoiceListResponse struct {
	Invoices   []*Invoice `json:"invoices"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

type InvoiceStats struct {
	TotalInvoices   int     `json:"total_invoices"`
	DraftInvoices   int     `json:"draft_invoices"`
	SentInvoices    int     `json:"sent_invoices"`
	PaidInvoices    int     `json:"paid_invoices"`
	OverdueInvoices int     `json:"overdue_invoices"`
	TotalRevenue    float64 `json:"total_revenue"`
	PendingRevenue  float64 `json:"pending_revenue"`
	OverdueRevenue  float64 `json:"overdue_revenue"`
}

// some constants related to invoice status

const (
	InvoiceStatusDraft    = "draft"
	InvoiceStatusSent     = "sent"
	InvoiceStatusPaid     = "paid"
	InvoiceStatusOverdue  = "overdue"
	InvoiceStatusCanceled = "canceled"
)

// some template constants

const (
	TemplateDefault      = "default"
	TemplateModern       = "modern"
	TemplateMinimal      = "minimal"
	TemplateProfessional = "professional"
)
