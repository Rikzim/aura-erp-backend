-- ── Database: crud_db ─────────────────────────────────────────────────────────

-- ── Users ─────────────────────────────────────────────────────────────────────
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'sales', 'production', 'viewer')),
    created_at      TIMESTAMP DEFAULT NOW()
);

-- ── Clients ───────────────────────────────────────────────────────────────────
CREATE TABLE clients (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(100),
    phone           VARCHAR(30),
    address         TEXT,
    vat_number      VARCHAR(30),
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- ── Sections ──────────────────────────────────────────────────────────────────
CREATE TABLE sections (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- ── Products ──────────────────────────────────────────────────────────────────
CREATE TABLE products (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    unit_price      NUMERIC(10, 2) NOT NULL,
    unit            VARCHAR(20) NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- ── Proposals ─────────────────────────────────────────────────────────────────
CREATE TABLE proposals (
    id              SERIAL PRIMARY KEY,
    reference       VARCHAR(50) NOT NULL UNIQUE,
    client_id       INTEGER REFERENCES clients(id) ON DELETE RESTRICT,
    section_id      INTEGER REFERENCES sections(id) ON DELETE SET NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'pending', 'approved', 'rejected')),
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- ── Proposal Items ────────────────────────────────────────────────────────────
CREATE TABLE proposal_items (
    id              SERIAL PRIMARY KEY,
    proposal_id     INTEGER NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
    product_id      INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity        NUMERIC(10, 2) NOT NULL,
    unit_price      NUMERIC(10, 2) NOT NULL,
    notes           TEXT
);

-- ── Orders ────────────────────────────────────────────────────────────────────
CREATE TABLE orders (
    id              SERIAL PRIMARY KEY,
    reference       VARCHAR(50) NOT NULL UNIQUE,
    proposal_id     INTEGER REFERENCES proposals(id) ON DELETE SET NULL,
    client_id       INTEGER REFERENCES clients(id) ON DELETE RESTRICT,
    section_id      INTEGER REFERENCES sections(id) ON DELETE SET NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_production', 'completed', 'delivered')),
    due_date        DATE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- ── Order Items ───────────────────────────────────────────────────────────────
CREATE TABLE order_items (
    id              SERIAL PRIMARY KEY,
    order_id        INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id      INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity        NUMERIC(10, 2) NOT NULL,
    unit_price      NUMERIC(10, 2) NOT NULL,
    notes           TEXT
);

-- ── Audit Log ─────────────────────────────────────────────────────────────────
CREATE TABLE audit_log (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       INTEGER NOT NULL,
    action          VARCHAR(50) NOT NULL,
    old_value       JSONB,
    new_value       JSONB,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- ── Seed Data ─────────────────────────────────────────────────────────────────

-- Users (passwords are SHA256 hashed)
-- Main admin: email=admin@aura-erp.com, password=admin123
-- Other users: password=changeme
INSERT INTO users (name, email, password_hash, role) VALUES
    ('Admin User',     'admin@aura-erp.com',     '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9', 'admin'),
    ('Admin User',     'admin@company.com',      '5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5', 'admin'),
    ('Sales User',     'sales@company.com',      '5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5', 'sales'),
    ('Production User','production@company.com', '5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5', 'production');


INSERT INTO clients (name, email, phone, address, vat_number) VALUES
    ('ACME Corp',   'contact@acme.com',   '+351 910 000 001', 'Rua de Lisboa 1, Porto',   'PT123456789'),
    ('Globex Corp', 'contact@globex.com', '+351 910 000 002', 'Avenida Central 5, Lisboa','PT987654321');

INSERT INTO sections (name, description) VALUES
    ('Metalwork',  'Metal cutting and welding'),
    ('Painting',   'Surface treatment and painting'),
    ('Assembly',   'Final assembly and QC');

INSERT INTO products (name, description, unit_price, unit) VALUES
    ('Steel Sheet 1mm',  '1mm cold rolled steel sheet', 12.50, 'm2'),
    ('Steel Tube 40x40', '40x40 square steel tube',     8.00,  'm'),
    ('Paint Primer',     'Grey epoxy primer',            5.00,  'kg');

INSERT INTO proposals (reference, client_id, section_id, status, notes) VALUES
    ('PROP-2024-001', 1, 1, 'draft',   'First proposal for ACME'),
    ('PROP-2024-002', 2, 2, 'pending', 'Awaiting client approval');

INSERT INTO proposal_items (proposal_id, product_id, quantity, unit_price) VALUES
    (1, 1, 10.00, 12.50),
    (1, 2, 5.00,  8.00),
    (2, 3, 20.00, 5.00);

INSERT INTO orders (reference, proposal_id, client_id, section_id, status, due_date) VALUES
    ('ORD-2024-001', 1, 1, 1, 'in_production', '2024-03-01');

INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
    (1, 1, 10.00, 12.50),
    (1, 2, 5.00,  8.00);

INSERT INTO audit_log (user_id, entity_type, entity_id, action, old_value, new_value) VALUES
    (1, 'proposal', 1, 'created',       NULL,                          '{"status": "draft"}'),
    (2, 'proposal', 1, 'status_changed','{"status": "draft"}',         '{"status": "pending"}'),
    (1, 'proposal', 1, 'converted',     '{"status": "approved"}',      '{"order_ref": "ORD-2024-001"}');
