//Core business logic user-service

package service

import (
	"time"

	"github.com/Suthar345Piyush/invoicego/internal/database"
	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/Suthar345Piyush/invoicego/internal/util"
	"github.com/google/uuid"
)

// structs for user service related work
type UserService struct {
	db *database.DB
}

// struct for new user service

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

// creating user

func (s *UserService) CreateUser(req *domain.RegisterRequest) (*domain.User, error) {

	// checking if user already exists or not
	// executing an query to get atleast one row

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)

	if err != nil {
		return nil, err
	}

	//  if user already exists , then , just returning user already exists

	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// hashing the password

	hashsedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// creating the user

	user := &domain.User{

		ID:                  uuid.New(),
		Email:               req.Email,
		PasswordHash:        hashsedPassword,
		FullName:            req.FullName,
		SubscriptionTier:    "free",
		SubscriptionStatus:  "active",
		MonthlyInvoiceCount: 0,
		MonthlyInvoiceLimit: 5,
		DefaultCurrency:     "INR",
		DefaultPaymentTerms: 30,
		InvoiceNumberPrefix: "INV",
		NextInvoiceNumber:   1,
		EmailVerified:       false,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// writing SQL queries

	query :=
		`
		       INSERT INTO users (
						   id , email , password_hash , full_name , subscription_tier , subscription_status , monthly_invoice_count , monthly_invoice_limit , default_currency , default_payment_terms , 
							 invoice_number_prefix , next_invoice_number , email_verified , is_active , created_at , updated_at
					 ) VALUES ($1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10 , $11 , $12 , $13 , $14 , $15 , $16)
		   `

	// executing query without returning

	_, err = s.db.Exec(
		query,
		user.ID, user.Email, user.PasswordHash, user.FullName, user.SubscriptionTier, user.SubscriptionStatus, user.MonthlyInvoiceCount, user.MonthlyInvoiceLimit, user.DefaultCurrency,
		user.DefaultPaymentTerms, user.InvoiceNumberPrefix, user.NextInvoiceNumber, user.EmailVerified, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}

// after creating user , we getting user by their email and their id

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {

	user := &domain.User{}

	query :=
		`
			    SELECT id , email , password_hash , full_name , business_name , business_address , business_phone , business_email , tax_id , logo_url , subscription_tier , subscription_status , monthly_invoice_count , monthly_invoice_limit , default_currency , default_payment_terms , invoice_number_prefix  , next_invoice_number , email_verified , is_active , created_at , updated_at
					last_login_at FROM users WHERE email = $1 AND is_active = true   
			  `

	// queryRow at most returns a row after querying the table

	err := s.db.QueryRow(query, email).Scan()

}
