-- Users table and it's components 

CREATE EXTENSION IF NOT EXISTS  "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     email VARCHAR(255) UNIQUE NOT NULL,
     password_hash VARCHAR(255) NOT NULL,
     full_name VARCHAR(255) NOT NULL,

     

     -- BUSINESS INFORMATION 

     business_name VARCHAR(255),
     business_address TEXT,
     business_phone VARCHAR(50),
     business_email VARCHAR(255),
     tax_id VARCHAR(100),
     logo_url TEXT,



     --subscription and limits 

     subscription_tier VARCHAR(50) DEFAULT 'free',
     subscription_status VARCHAR(50) DEFAULT 'active',
     subscription_started_at TIMESTAMP,
     subscription_expires_at TIMESTAMP,
     monthly_invoice_count INT DEFAULT 0,
     monthly_invoice_limit INT DEFAULT 5,


    default_currency VARCHAR(3) DEFAULT 'INR',
    default_payment_terms INT DEFAULT 30,
    invoice_number_prefix VARCHAR(20) DEFAULT 'INV',
    next_invoice_number INT DEFAULT 1,



    email_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP

);



-- creating indexes on user's email and subscription 

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_subscription_tier ON users(subscription_tier);
