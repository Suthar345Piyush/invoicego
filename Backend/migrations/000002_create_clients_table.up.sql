-- postgres migrations for clients 

CREATE TABLE IF NOT EXISTS clients (

   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,


   --Client info--

   name VARCHAR(255) NOT NULL,
   email VARCHAR(255),
   phone VARCHAR(50),
   company_name VARCHAR(255),



   --client address info--

   address_line1 VARCHAR(255),
   address_line2 VARCHAR(255),
   city VARCHAR(100),
   state VARCHAR(100),
   postal_code VARCHAR(20),
   country VARCHAR(100),

   tax_id VARCHAR(100),
   notes TEXT,

   is_active BOOLEAN DEFAULT true,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

  -- indexes on user id , email and name 

  CREATE INDEX idx_clients_user_id ON clients(user_id);
  CREATE INDEX idx_clients_name ON clients(name);
  CREATE INDEX idx_clients_email ON clients(email);
