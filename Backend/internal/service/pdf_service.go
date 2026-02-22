// pdf service code here

package service

import (
	"fmt"
	"time"

	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

// pdf service function

func NewPDFService() *PDFService {
	return &PDFService{}
}

//  invoice pdf function

func (s *PDFService) GenerateInvoicePDF(invoice *domain.Invoice, user *domain.User) ([]byte, error) {

	// writing pdf conventions

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// setting font

	pdf.SetFont("Arial", "", 12)

	// header of the invoice - business information

	s.addHeader(pdf, user, invoice)

	// client information

	s.addClientInfo(pdf, invoice.Client)

	// invoice details

	s.addInvoiceDetails(pdf, invoice)

	// line items table

	s.addItemsTable(pdf, invoice.Items)

	// totals

	s.addTotals(pdf, invoice)

	// notes and terms

	s.addNotesAndTerms(pdf, invoice)

	// footer part of the invoice

	s.addFooter(pdf)

	//getting the pdf as bytes

	var buf []byte
	buf, err := pdf.Output(New.Addoimcde)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// add header function

func (s *PDFService) addHeader(pdf *gofpdf.Fpdf, user *domain.User, invoice *domain.Invoice) {

	// company name

	pdf.SetFont("Arial", "B", 20)

	if user.BusinessName != nil && *user.BusinessName != "" {
		pdf.Cell(0, 10, *user.BusinessName)
	} else {
		pdf.Cell(0, 10, user.FullName)
	}

	pdf.Ln(8)

	// company details

	pdf.SetFont("Arial", "", 10)

	if user.BusinessAddress != nil && *user.BusinessAddress != "" {
		pdf.Cell(0, 5, *user.BusinessAddress)
		pdf.Ln(5)
	}

	if user.BusinessEmail != nil && *user.BusinessEmail != "" {
		pdf.Cell(0, 5, "Email : "+*user.BusinessEmail)
		pdf.Ln(5)
	}

	if user.BusinessPhone != nil && *user.BusinessPhone != "" {
		pdf.Cell(0, 5, "Phone: "+*user.BusinessPhone)
		pdf.Ln(5)
	}

	if user.TaxID != nil && *user.TaxID != "" {
		pdf.Cell(0, 5, "Tax ID "+*user.TaxID)
		pdf.Ln(5)
	}

	pdf.Ln(10)

	//invoice title section

	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(0, 102, 204)
	pdf.Cell(0, 10, "INVOICE")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(12)

}

//  client information function

func (s *PDFService) addClientInfo(pdf *gofpdf.Fpdf, client *domain.Client) {

	if client == nil {
		return
	}

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, "Bill To:")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, client.Name)
	pdf.Ln(5)

	if client.CompanyName != nil && *client.CompanyName != "" {
		pdf.Cell(0, 5, *client.CompanyName)
		pdf.Ln(5)
	}

	if client.AddressLine1 != nil && *client.AddressLine1 != "" {
		pdf.Cell(0, 5, *client.AddressLine1)
		pdf.Ln(5)
	}

	if client.AddressLine2 != nil && *client.AddressLine2 != "" {
		pdf.Cell(0, 5, *client.AddressLine2)
		pdf.Ln(5)
	}

	// address parts like state , city , Postal code

	addressParts := ""

	if client.City != nil && *client.City != "" {
		addressParts += *client.City
	}

	if client.State != nil && *client.State != "" {
		if addressParts != "" {
			addressParts += ", "
		}
		addressParts += *client.State
	}

	if client.PostalCode != nil && *client.PostalCode != "" {
		if addressParts != "" {
			addressParts += " "
		}
		addressParts += *client.PostalCode
	}

	if addressParts != "" {
		pdf.Cell(0, 5, addressParts)
		pdf.Ln(5)
	}

	if client.Email != nil && *client.Email != "" {
		pdf.Cell(0, 5, "Email: "+*client.Email)
		// ln - line break
		pdf.Ln(5)
	}

	pdf.Ln(8)

}

// invoice details  function like issue date , invoice number , status , due date

func (s *PDFService) addInvoiceDetails(pdf *gofpdf.Fpdf, invoice *domain.Invoice) {

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Invoice Number:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, invoice.InvoiceNumber)
	pdf.Ln(6)

	// issue date
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Issue Date:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, invoice.IssueDate.Format("February 19, 2026"))
	pdf.Ln(6)

	//due date
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Due Date:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, invoice.DueDate.Format("February 19, 2026"))
	pdf.Ln(6)

	// status
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 6, "Status:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, invoice.Status)
	pdf.Ln(10)

}

//  adding items table

func (s *PDFService) addItemsTable(pdf *gofpdf.Fpdf, items []*domain.InvoiceItem) {

	// table header

	pdf.SetFillColor(200, 220, 255)
	pdf.SetFont("Arial", "B", 10)

	pdf.CellFormat(90, 8, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, "Unit Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(35, 8, "Amount", "1", 1, "R", true, 0, "")

	// table body

	pdf.SetFont("Arial", "", 10)

	for _, item := range items {

		pdf.CellFormat(90, 7, item.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 7, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 7, fmt.Sprintf("%.2f", item.Amount), "1", 1, "R", false, 0, "")

	}
	pdf.Ln(5)
}

// add totals function

func (s *PDFService) addTotals(pdf *gofpdf.Fpdf, invoice *domain.Invoice) {

	// positions for the totals

	startX := 120.0

	pdf.SetX(startX)
	pdf.Cell(35, 6, "Subtotal:")
	pdf.Cell(35, 6, fmt.Sprintf("%s %.2f", invoice.Currency, invoice.Subtotal))
	pdf.Ln(6)

	// tax part of the invoice

	if invoice.TaxRate > 0 {
		pdf.SetX(startX)
		pdf.Cell(35, 6, fmt.Sprintf("Tax (%.2f%%):", invoice.TaxRate))
		pdf.Cell(35, 6, fmt.Sprintf("%s %.2f", invoice.Currency, invoice.TaxAmount))
		pdf.Ln(6)
	}

	// discount is greater than zero

	if invoice.DiscountAmount > 0 {
		pdf.SetX(startX)
		pdf.Cell(35, 6, "Discount:")
		pdf.Cell(35, 6, fmt.Sprintf("-%s %.2f", invoice.Currency, invoice.DiscountAmount))
		pdf.Ln(6)
	}

	// final - total

	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(startX)
	pdf.Cell(35, 8, "Total:")
	pdf.Cell(35, 8, fmt.Sprintf("%s %.2f", invoice.Currency, invoice.TotalAmount))
	pdf.Ln(12)

}

// add notes and terms of the invoice

func (s *PDFService) addNotesAndTerms(pdf *gofpdf.Fpdf, invoice *domain.Invoice) {

	if invoice.Notes != nil && *invoice.Notes != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, "Notes:")
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 9)
		pdf.MultiCell(0, 5, *invoice.Notes, "", "", false)
		pdf.Ln(5)
	}

	if invoice.TermsAndConditions != nil && *invoice.TermsAndConditions != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, "Terms & Conditions:")
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 9)
		pdf.MultiCell(0, 5, *invoice.TermsAndConditions, "", "", false)
		pdf.Ln(5)
	}

}

// footer part of the invoice

func (s *PDFService) addFooter(pdf *gofpdf.Fpdf) {
	pdf.SetY(-20)
	pdf.SetFont("Arial", "1", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", time.Now().Format("February 19 2026")))
}
