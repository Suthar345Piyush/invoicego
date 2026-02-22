-- QUERIES FOR THE INVOICES TABLE  

CREATE TABLE IF NOT EXISTS  invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,



    -- invoice details query  

    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',



    -- TYPE OF DATES 

    issue_date DATE NOT NULL,
    due_date  DATE NOT NULL,
    paid_date DATE,


    -- amount and currency related things 

    currency VARCHAR(3) DEFAULT 'INR',
    subtotal DECIMAL(15 , 2) NOT NULL DEFAULT 0,
    tax_rate DECIMAL(5 , 2) DEFAULT 0,
    tax_amount DECIMAL(15 , 2) DEFAULT 0,
    discount_amount  DECIMAL(15 , 2) DEFAULT 0,
    total_amount DECIMAL(15 , 2) NOT NULL DEFAULT 0,


    -- TEMPLATE CUSTOMIZATION 

    template_id VARCHAR(50) DEFAULT 'default',
    notes TEXT,
    terms_and_conditions TEXT,


    --- pdf url , time it was generated 

    pdf_url TEXT,
    pdf_generated_at TIMESTAMP,


    -- email related things 

    email_sent BOOLEAN DEFAULT false,
    email_sent_at TIMESTAMP,
    email_opened BOOLEAN DEFAULT false,
    email_opened_at TIMESTAMP,


    -- created_At and updated_At 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);



--- creating invoices_items table  


CREATE TABLE IF NOT EXISTS invoice_items (
    
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,


     -- items details  

     description TEXT NOT NULL,
     quantity DECIMAL(10 , 2) NOT NULL DEFAULT 1,
     unit_price DECIMAL(15 , 2) NOT NULL,
     amount DECIMAL(15 , 2) NOT NULL,


     -- order  
     sort_order INT NOT NULL DEFAULT 0,

     
     -- created_at , updated_at 

     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);


-- indexes on invoices  

CREATE INDEX idx_invoices_user_id ON invoices(user_id);
CREATE INDEX idx_invoices_client_id ON invoices(client_id);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
CREATE INDEX idx_invoices_issue_date ON invoices(issue_date);



-- indexes in invoice items tables 


CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);
CREATE INDEX idx_invoice_items_sort_order ON invoice_items(invoice_id , sort_order);


