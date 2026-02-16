// client file in domain

package domain

import (
	"time"

	"github.com/google/uuid"
)

// client information struct
// structs with the tag value
type Client struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Email        *string   `json:"email,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	CompanyName  *string   `json:"company_name,omitempty"`
	AddressLine1 *string   `json:"address_line1,omitempty"`
	AddressLine2 *string   `json:"address_line2,omitempty"`
	City         *string   `json:"city,omitempty"`
	State        *string   `json:"state,omitempty"`
	PostalCode   *string   `json:"postal_code,omitempty"`
	Country      *string   `json:"country,omitempty"`
	TaxID        *string   `json:"tax_id,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// creating client request struct

type CreateClientRequest struct {
	Name         string  `json:"name" validate:"required,min=2"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone        *string `json:"phone,omitempty"`
	CompanyName  *string `json:"company_name,omitempty"`
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      *string `json:"country,omitempty"`
	TaxID        *string `json:"tax_id,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

// update client request struct

type UpdateClientRequest struct {
	Name         string  `json:"name" validate:"required,min=2"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone        *string `json:"phone,omitempty"`
	CompanyName  *string `json:"company_name,omitempty"`
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      *string `json:"country,omitempty"`
	TaxID        *string `json:"tax_id,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

//client response list

type ClientListResponse struct {
	Clients   []*Client `json:"clients"`
	Total     int       `json:"total"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
	TotalPage int       `json:"total_page"`
}
