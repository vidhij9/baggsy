CREATE TABLE bags (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    bag_type VARCHAR(50) NOT NULL,
    status VARCHAR(50)
);
