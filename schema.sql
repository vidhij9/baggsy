CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('employee', 'admin'))
);

CREATE TABLE bags (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    type ENUM('parent', 'child') NOT NULL,
    child_count INT DEFAULT 0 CHECK (child_count >= 0),
    parent_id INT,
    linked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES bags(id),
    CONSTRAINT check_parent_child CHECK ((type = 'parent' AND parent_id IS NULL) OR (type = 'child' AND parent_id IS NOT NULL))
);
CREATE INDEX idx_bags_qr_code ON bags(lower(qr_code));
CREATE INDEX idx_bags_parent_id ON bags(parent_id);

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    parent_id INT UNIQUE NOT NULL,
    bill_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_link_parent FOREIGN KEY (parent_id) REFERENCES bags(id)
);
CREATE INDEX idx_links_bill_id ON links(bill_id);