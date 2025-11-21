CREATE TYPE application_type AS ENUM ('AUTO', 'PERSONAL');        
CREATE TYPE application_status AS ENUM ('NEW', 'REVIEW', 'APPROVED', 'REJECTED');  -- status: NEW / REVIEW / APPROVED / REJECTED
CREATE TYPE loan_status AS ENUM ('ACTIVE', 'PAID', 'OVERDUE');  -- status: ACTIVE / PAID / OVERDUE

CREATE TABLE IF NOT EXISTS loan_applications  (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT NOT NULL,
    type             application_type NOT NULL,
    vehicle_vin      VARCHAR(32),
    vehicle_name     VARCHAR(255),
    currency_code    VARCHAR(10) NOT NULL,
    price            NUMERIC(18,2),
    down_payment     NUMERIC(18,2),
    net_price        NUMERIC(18,2),
    margin_rate      NUMERIC(5,2),
    term_months      INT,
    monthly_payment  NUMERIC(18,2),
    status           application_status DEFAULT 'NEW', 
    created_at       TIMESTAMP DEFAULT NOW(),
    updated_at       TIMESTAMP
);

CREATE TABLE IF NOT EXISTS loans (
    id                 BIGSERIAL PRIMARY KEY,
    application_id     BIGINT REFERENCES loan_applications(id) NOT NULL,
    user_id            BIGINT NOT NULL,
    vehicle_vin        VARCHAR(32),
    currency_code      VARCHAR(10) NOT NULL,
    amount             NUMERIC(18,2),
    term_months        INT,
    monthly_payment    NUMERIC(18,2),
    remaining_balance  NUMERIC(18,2),
    status             loan_status DEFAULT 'ACTIVE',   
    created_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_loan_applications_user ON loan_applications(user_id);
CREATE INDEX idx_loans_user ON loans(user_id);
