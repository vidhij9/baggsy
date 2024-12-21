CREATE TABLE bags (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    bag_type VARCHAR(50) NOT NULL,
    status VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    sap_bill_id VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

ALTER TABLE bags ADD COLUMN bill_id INT REFERENCES bills(id);
