package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Suthar345Piyush/invoicego/internal/database"
	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/google/uuid"
)

type ClientService struct {
	db *database.DB
}

func NewClientService(db *database.DB) *ClientService {
	return &ClientService{db: db}
}

func (s *ClientService) CreateClient(userID uuid.UUID, req *domain.CreateClientRequest) (*domain.Client, error) {

	client := &domain.Client{

		ID:           uuid.New(),
		UserID:       userID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		CompanyName:  req.CompanyName,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		TaxID:        req.TaxID,
		Notes:        req.Notes,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query :=
		`
		       INSERT INTO clients (
						 id , user_id , name , email , phone , company_name , address_line1 , address_line2 , city , state , postal_code , country , tax_id , notes , is_active , created_at , updated_at
					 )  VALUES ($1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10 , $11 , $12 , $13 , $14 , $15, $16 , $17)
		   `

	_, err := s.db.Exec(
		query,
		client.ID, client.UserID, client.Name, client.Email, client.Phone, client.CompanyName, client.AddressLine1, client.AddressLine2, client.City, client.State, client.PostalCode, client.Country, client.TaxID, client.Notes, client.IsActive, client.CreatedAt, client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return client, nil

}

// getting client by  their client id's

func (s *ClientService) GetClientByID(userID, clientID uuid.UUID) (*domain.Client, error) {
	client := &domain.Client{}

	query :=
		`
		  SELECT id , user_id , name , email , phone , company_name , address_line1 , address_line2 , city , state,
			postal_code , tax_id , notes , is_active , created_at , updated_at FROM clients WHERE id = $1 AND user_id = $2 AND is_active = true
		 `

	err := s.db.QueryRow(query, clientID, userID).Scan(
		&client.ID, &client.UserID, &client.Name, &client.Email, &client.CompanyName, &client.AddressLine1, &client.AddressLine2, &client.City, &client.State, &client.PostalCode, &client.TaxID, &client.Notes, &client.IsActive, &client.CreatedAt, &client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("client not found")
	}

	if err != nil {
		return nil, err
	}

	return client, nil

}

// getting user by their user id , response in return

func (s *ClientService) GetClientsByUserID(userID uuid.UUID, page, pageSize int) (*domain.ClientListResponse, error) {

	// default pagination values

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	//getting total count

	var total int

	countQuery := `SELECT COUNT(*) FROM clients WHERE user_id = $1 AND is_Active = true`
	err := s.db.QueryRow(countQuery, userID).Scan(&total)

	if err != nil {
		return nil, err
	}

	// query for  getting clients

	query :=
		` SELECT id , user_id , name , email , phone , company_name , address_line1 , address_line2 ,
				 city , state , postal_code , country  , tax_id , notes , is_active , created_at , updated_at FROM clients WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC LIMIT $2 offset $3 
			  `

	rows, err := s.db.Query(query, userID, pageSize, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clients := []*domain.Client{}

	for rows.Next() {
		client := &domain.Client{}

		err := rows.Scan(
			&client.ID, &client.UserID, &client.Name, &client.Email, &client.Phone, &client.CompanyName, &client.AddressLine1, &client.AddressLine2, &client.City, &client.State, &client.PostalCode, &client.Country, &client.TaxID, &client.Notes, &client.IsActive, &client.CreatedAt, &client.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		clients = append(clients, client)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	//  total pages

	totalPages := (total + pageSize - 1) / pageSize

	return &domain.ClientListResponse{
		Clients:   clients,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPages,
	}, nil

}

// updating the client , it will return updated client

func (s *ClientService) UpdateClient(userID, clientID uuid.UUID, req *domain.UpdateClientRequest) (*domain.Client, error) {

	// checking if client exists and belongs to user or not

	_, err := s.GetClientByID(userID, clientID)

	if err != nil {
		return nil, err
	}

	// query to update the client dynamically
	// using COALESCE function for returning first non-null argument

	query := `UPDATE clients SET
			            name = COALESCE($1 , name),
									email = COALESCE($2 , email),
									phone = COALESCE($3 , phone),
									company_name = COALESCE($4 , company_name),
									address_line1 = COALESCE($5 , address_line1),
									address_line2 = COALESCE($6 , address_line2),
									city = COALESCE($7 , city),
									state = COALESCE($8 , state),
									postal_code = COALESCE($9 , postal_code),
									country = COALESCE($10 , country),
									tax_id = COALESCE($11 , tax_id),
									notes = COALESCE($12 , notes),
									updated_at = $13

								WHERE id = $14 AND  user_id = $15
			        `

	_, err = s.db.Exec(
		query,
		req.Name, req.Email, req.Phone, req.CompanyName, req.AddressLine1, req.AddressLine2, req.City, req.State, req.PostalCode, req.Country, req.TaxID, req.Notes, time.Now(), clientID, userID,
	)

	if err != nil {
		return nil, err
	}

	// return updated client

	return s.GetClientByID(userID, clientID)

}

// deleting the client

func (s *ClientService) DeleteClient(userID, clientID uuid.UUID) error {

	// checking if client exists and belongs to user or not

	_, err := s.GetClientByID(userID, clientID)

	if err != nil {
		return err
	}

	//deleting client carefully

	query :=
		`UPDATE clients SET is_active = false , updated_at = $1 WHERE id = $2 AND user_id = $3`

	_, err = s.db.Exec(query, time.Now(), clientID, userID)

	return err

}
