// invoice service

package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Suthar345Piyush/invoicego/internal/database"
	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/google/uuid"
)

type InvoiceService struct {
	db          *database.DB
	userService *UserService
}

// invoice service function

func NewInvoiceService(db *database.DB, userService *UserService) *InvoiceService {
	return &InvoiceService{
		db:          db,
		userService: userService,
	}
}

// create invoice function

func (s *InvoiceService) CreateInvoice(userID uuid.UUID, req *domain.CreateInvoiceRequest) (*domain.Invoice, error) {

	// getting user to checking his subscription limits

	user, err := s.userService.GetUserByID(userID)

	if err != nil {
		return nil, err
	}

	// initially free tier limit check

	if user.SubscriptionTier == "free" && user.MonthlyInvoiceCount >= user.MonthlyInvoiceLimit {
		return nil, domain.ErrInvoiceLimitExceeded
	}

	// parsing dates

	issueDate, err := time.Parse("2026-02-20", req.IssueDate)

	if err != nil {
		return nil, fmt.Errorf("invalid issue_date format, use YYYY-MM-DD")
	}

	dueDate, err := time.Parse("2026-02-20", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format , user YYYY-MM-DD")
	}

	// calculatin of amounts

	subtotal := 0.0

	for _, item := range req.Items {
		subtotal += item.Quantity * item.UnitPrice
	}

	taxAmount := (subtotal * req.TaxRate) / 100
	totalAmount := subtotal + taxAmount - req.DiscountAmount

	// generating invoice number

	invoiceNumber := fmt.Sprintf("%s-%04d", user.InvoiceNumberPrefix, user.NextInvoiceNumber)

	// setting defualt template
	templateID := req.TemplateID

	if templateID == "" {
		templateID = domain.TemplateDefault
	}

	// creating invoice in key : value format

	invoice := &domain.Invoice{

		ID:                 uuid.New(),
		UserID:             userID,
		ClientID:           req.ClientID,
		InvoiceNumber:      invoiceNumber,
		Status:             domain.InvoiceStatusDraft,
		IssueDate:          issueDate,
		DueDate:            dueDate,
		Currency:           req.Currency,
		Subtotal:           subtotal,
		TaxRate:            req.TaxRate,
		TaxAmount:          taxAmount,
		DiscountAmount:     req.DiscountAmount,
		TotalAmount:        totalAmount,
		TemplateID:         templateID,
		Notes:              req.Notes,
		TermsAndConditions: req.TermsAndConditions,
		EmailSent:          false,
		EmailOpened:        false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// starting initial transaction

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	// aborting the transaction before , complete function execution

	defer tx.Rollback()

	// insert into invoices table

	invoiceQuery :=

		`INSERT INTO invoices (
			   id , user_id , client_id , invoice_number , status , issue_date , due_date , currency , subtotal , tax_rate , tax_amount , discount_amount , total_amount , template_id , notes , terms_and_conditions , email_sent , email_opened , created_at , updated_at   
		 ) VALUES ($1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10 , $11 , $12 , $12 , $13 , $14 , $15 , $16 , $17 , $18 , $19 , $20)`

	_, err = tx.Exec(
		invoiceQuery,
		invoice.ID, invoice.UserID, invoice.ClientID, invoice.InvoiceNumber, invoice.Status, invoice.IssueDate, invoice.DueDate, invoice.Currency, invoice.Subtotal, invoice.TaxRate, invoice.TaxAmount, invoice.DiscountAmount, invoice.TotalAmount, invoice.TemplateID, invoice.Notes, invoice.TermsAndConditions, invoice.EmailSent, invoice.EmailOpened, invoice.CreatedAt, invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// inserting invoice items into table

	itemQuery := `
		      INSERT INTO invoice_items (
					  	 id , invoice_id , description , quantity , unit_price , amount , sort_order , created_at , updated_at
					) VALUES ($1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9)
		    `

	// creating item id and  amount = item(quantity * unitPrice)

	for i, item := range req.Items {

		itemID := uuid.New()
		amount := item.Quantity * item.UnitPrice

		_, err = tx.Exec(
			itemQuery,
			itemID, invoice.ID, item.Description, item.Quantity, item.UnitPrice, amount, i, time.Now(), time.Now(),
		)

		if err != nil {
			return nil, err
		}

	}

	// updating users table  next invoice number and  monthly count  , both by one

	updateUserQuery :=

		` UPDATE users SET 
					  next_invoice_number  = next_invoice_number + 1,
						monthly_invoice_count = monthly_invoice_count + 1,
						updated_at = $1 WHERE id = $2
				  `

	_, err = tx.Exec(updateUserQuery, time.Now(), userID)

	if err != nil {
		return nil, err
	}

	// commiting transaction

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	//  getting full invoice with the items and client

	return s.GetInvoiceByID(userID, invoice.ID)

}

// getting full invoice by id

func (s *InvoiceService) GetInvoiceByID(userID, invoiceID uuid.UUID) (*domain.Invoice, error) {

	invoice := &domain.Invoice{}

	// query on invoices table  with id and user id

	query :=
		`SELECT id , user_id , client_id , invoice_number , status , issue_date , due_date , currency , subtotal , tax_rate , tax_amount , dicount_amount , total_amount , template_id , notes , terms_and_conditions , pdf_url , pdf_generated_at , email_sent , email_sent_at , email_opened , email_opened_at , created_at , updated_at  FROM invoices WHERE id = $1 AND user_id = $2`

	err := s.db.QueryRow(query, invoiceID, userID).Scan(
		&invoice.ID, &invoice.UserID, &invoice.ClientID, &invoice.InvoiceNumber, &invoice.Status, &invoice.IssueDate, &invoice.DueDate, &invoice.Currency, &invoice.Subtotal, &invoice.TaxRate, &invoice.TaxAmount, &invoice.DiscountAmount, &invoice.TotalAmount, &invoice.TemplateID, &invoice.Notes, &invoice.TermsAndConditions, &invoice.PDFURL, &invoice.PDFGeneratedAt, &invoice.EmailSent, &invoice.EmailSentAt, &invoice.EmailOpened, &invoice.EmailOpenedAt, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found")
	}

	if err != nil {
		return nil, err
	}

	// getting invoice items

	items, err := s.getInvoiceItems(invoiceID)

	if err != nil {
		return nil, err
	}

	invoice.Items = items

	// getting client information

	clientQuery :=
		`
			     SELECT id , user_id , name , email , phone , comapny_name , address_line1 , address_line2 , city , state , postal_code , country , tax_id , notes , is_active , created_at , updated_at FROM clients WHERE id = $1 
			   `

	// getting client

	client := &domain.Client{}

	err = s.db.QueryRow(clientQuery, invoice.ClientID).Scan(
		&client.ID, &client.UserID, &client.Name, &client.Email, &client.Phone, &client.CompanyName, &client.AddressLine1, &client.AddressLine2, &client.City, &client.State, &client.PostalCode, &client.Country, &client.TaxID, &client.Notes, &client.IsActive, &client.CreatedAt, &client.UpdatedAt,
	)

	if err == nil {
		invoice.Client = client
	}

	return invoice, nil

}

// function for getting invoice items
// returns an array of invoice items

func (s *InvoiceService) getInvoiceItems(invoiceID uuid.UUID) ([]*domain.InvoiceItem, error) {

	query :=
		`
		     SELECT id , invoice_id , description , quantity , unit_price , amount , sort_order , created_at , updated_at FROM invoice_items WHERE invoice_id = $1 ORDER BY sort_order 
		   `

	rows, err := s.db.Query(query, invoiceID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []*domain.InvoiceItem{}

	for rows.Next() {
		item := &domain.InvoiceItem{}

		err := rows.Scan(
			&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Amount, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

// function for  getting invoices by the user id
// list of invoices are returned

func (s *InvoiceService) GetInvoiceByUserID(userID uuid.UUID, page, pageSize int, status string) (*domain.InvoiceListResponse, error) {

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// building query with optional status filter

	countQuery := `SELECT COUNT(*) FROM invoices WHERE user_id = $1`

	args := []interface{}{userID}

	if status != "" {
		countQuery += ` AND status = %2`
	args:
		append(args, status)
	}

	var total int

	err := s.db.QueryRow(countQuery, args...).Scan(&total)

	if err != nil {
		return nil, err
	}

	// query to get invoices

	query :=
		`
		      SELECT id , user_id , client_id , invoice_number , status , issue_date , due_date , currency , subtotal , tax_rate , tax_amount , dicount_amount , total_amount , template_id , notes , terms_and_conditions , pdf_url , pdf_generated_at , email_sent , email_sent_at , email_opened , email_opened_at , created_at , updated_at FROM invoices WHERE user_id = $1
		   `

	queryArgs := []interface{}{userID}

	if status != "" {
		query += ` AND status = $2`
		queryArgs = append(queryArgs, status)
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(queryArgs)+1) + `OFFSET $` + fmt.Sprintf("%d", len(queryArgs)+2)

	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := s.db.Query(query, queryArgs...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// getting invoices from domain in array format

	invoices := []*domain.Invoice{}

	for rows.Next() {

		invoice := &domain.Invoice{}

		err := rows.Scan(
			&invoice.ID, &invoice.UserID, &invoice.ClientID, &invoice.InvoiceNumber, &invoice.Status, &invoice.IssueDate, &invoice.DueDate, &invoice.Currency, &invoice.Subtotal, &invoice.TaxRate, &invoice.TaxAmount, &invoice.DiscountAmount, &invoice.TotalAmount, &invoice.TemplateID, &invoice.Notes, &invoice.TermsAndConditions, &invoice.PDFURL, &invoice.PDFGeneratedAt, &invoice.EmailSent, &invoice.EmailSentAt, &invoice.EmailOpened, &invoice.EmailOpenedAt, &invoice.CreatedAt, &invoice.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// calculating total pages

	totalPages := (total + pageSize - 1) / pageSize

	return &domain.InvoiceListResponse{
		Invoices:   invoices,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil

}

// function for updating invoice status

func (s *InvoiceService) UpdateInvoiceStatus(userID, invoiceID uuid.UUID, req *domain.UpdateInvoiceStatusRequest) (*domain.Invoice, error) {

	// getting that invoice to confirm , that user exists

	invoice, err := s.GetInvoiceByID(userID, invoiceID)

	if err != nil {
		return nil, err
	}

	//check - don't allow status update if invoice is already paid or canceled

	if invoice.Status == domain.InvoiceStatusPaid || invoice.Status == domain.InvoiceStatusCanceled {
		return nil, fmt.Errorf("cannot update status of %s invoice", invoice.Status)
	}

	query := `UPDATE invoices SET status = $1 , updated_at = $2`

	args := []interface{}{req.Status, time.Now()}

	argCount := 2

	// if invoice is paid , then setting paid_date

	if req.Status == domain.InvoiceStatusPaid {
		if req.PaidDate != nil {
			paidDate, err := time.Parse("2026-02-20", *req.PaidDate)

			if err != nil {
				return nil, fmt.Errorf("invalid paid_date format , use YYYY-MM-DD")
			}

			argCount++

			query += fmt.Sprintf(`, paid_date = $%d`, argCount)

			args = append(args, paidDate)

		} else {
			argCount++

			query += fmt.Sprintf(`, paid_date = $%d`, argCount)
			args = append(args, time.Now())

		}
	}

	argCount++

	query += fmt.Sprintf(` WHERE id = $%d AND user_id = $%d`, argCount, argCount+1)

	args = append(args, invoiceID, userID)

	// executing the query and arguments

	_, err = s.db.Exec(query, args...)

	if err != nil {
		return nil, err
	}

	return s.GetInvoiceByID(userID, invoiceID)

}

//  function to delete invoice

func (s *InvoiceService) DeleteInvoice(userID, invoiceID uuid.UUID) error {

	// delete only those invoices which are in draft

	invoice, err := s.GetInvoiceByID(userID, invoiceID)

	if err != nil {
		return err
	}

	if invoice.Status != domain.InvoiceStatusDraft {
		return fmt.Errorf("only draft invoices can be deleted")
	}

	// query to delete the invoice

	query := `DELETE FROM invoices WHERE id = $1 AND user_id = $2`

	// executing the query

	_, err = s.db.Exec(query, invoiceID, userID)

	return err

}

// function for invoices stats
// how much invoices are in which-which status (draft , paid , sent , total revenue , pending , overdue)

func (s *InvoiceService) GetInvoiceStats(userID uuid.UUID) (*domain.InvoiceStats, error) {

	stats := &domain.InvoiceStats{}

	// query to get stats of the invoices

	query := `
		     SELECT 
				    COUNT(*) as total_invoices
						COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_invoices,
						COUNT(CASE WHEN status = 'sent' THEN 1 END) as sent_invoices,
						COUNT(CASE WHEN status = 'paid' THEN 1 END) as paid_invoices,
						COUNT(CASE WHEN status = 'overdue' THEN 1 END) as overdue_invoices,
						COALESCE(SUM(CASE WHEN status = 'paid' THEN total_amount ELSE 0 END) , 0) as total_revenue,
						COALESCE(SUM(CASE WHEN status = 'sent' THEN total_amount ELSE 0 END) , 0) as pending_revenue,
						COALESCE(SUM(CASE WHEN status = 'overdue' THEN total_amount ELSE 0 END) , 0) as overdue_revenue

					FROM invoices
					WHERE user_id = $1
		  `

	err := s.db.QueryRow(query, userID).Scan(
		&stats.TotalInvoices, &stats.DraftInvoices, &stats.SentInvoices, &stats.PaidInvoices, &stats.OverdueInvoices, &stats.TotalRevenue, &stats.PendingRevenue, &stats.OverdueRevenue,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil

}

// function to create duplicate (copy) invoices

func (s *InvoiceService) DuplicateInvoice(userID, invoiceID uuid.UUID) (*domain.Invoice, error) {

	// to make copy we getting original invoice first

	originalInvoice, err := s.GetInvoiceByID(userID, invoiceID)

	if err != nil {
		return nil, err
	}

	// creating request from original invoice

	items := []*domain.CreateInvoiceItemReq{}

	for _, item := range originalInvoice.Items {
		items = append(items, &domain.CreateInvoiceItemReq{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}

	// creating invoice

	req := &domain.CreateInvoiceRequest{
		ClientID:           originalInvoice.ClientID,
		IssueDate:          time.Now().Format("2026-02-20"),
		DueDate:            time.Now().AddDate(0, 0, 30).Format("2026-02-20"),
		Currency:           originalInvoice.Currency,
		TaxRate:            originalInvoice.TaxRate,
		DiscountAmount:     originalInvoice.DiscountAmount,
		TemplateID:         originalInvoice.TemplateID,
		Notes:              originalInvoice.Notes,
		TermsAndConditions: originalInvoice.TermsAndConditions,
		Items:              items,
	}

	return s.CreateInvoice(userID, req)

}
