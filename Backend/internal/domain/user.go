// setting domain model - user
// user related types defination

package domain

import (
	"time"

	"github.com/google/uuid"
)

// setting some user things
// user's data returns in json format

type User struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"` // hash code of the password
	FullName            string     `json:"full_name"`
	BusinessName        *string    `json:"business_name,omitempty"`
	BusinessAddress     *string    `json:"business_address,omitempty"`
	BusinessPhone       *string    `json:"business_phone,omitempty"`
	BusinessEmail       *string    `json:"business_email,omitempty"`
	TaxID               *string    `json:"tax_id,omitempty"`
	LogoURL             *string    `json:"logo_url,omitempty"`
	SubscriptionTier    string     `json:"subscription_tier"`
	SubscriptionStatus  string     `json:"subscription_status"`
	MonthlyInvoiceCount int        `json:"monthly_invoice_count"`
	MonthlyInvoiceLimit int        `json:"monthly_invoice_limit"`
	DefaultCurrency     string     `json:"default_currency"`
	DefaultPaymentTerms int        `json:"default_payment_terms"`
	InvoiceNumberPrefix string     `json:"invoice_number_prefix"`
	NextInvoiceNumber   int        `json:"next_invoice_number"`
	EmailVerified       bool       `json:"email_verified"`
	IsActive            bool       `json:"is_active"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
}

// registration request struct

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
