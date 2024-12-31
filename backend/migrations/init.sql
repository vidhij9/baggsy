CREATE TABLE bags (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    bag_type VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bag_map (
    id SERIAL PRIMARY KEY,
    parent_bag VARCHAR(255) NOT NULL,
    child_bag VARCHAR(255) NOT NULL,
    UNIQUE (parent_bag, child_bag),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    parent_bag VARCHAR(255) NOT NULL,
    bill_id VARCHAR(255) NOT NULL,
    UNIQUE (parent_bag, bill_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bags_qr_code ON bags (qr_code);
CREATE INDEX idx_bag_map_parent_bag ON bag_map (parent_bag);
CREATE INDEX idx_bag_map_child_bag ON bag_map (child_bag);
CREATE INDEX idx_links_bill_id ON links (bill_id);

CREATE OR REPLACE FUNCTION soft_delete_bags()
RETURNS TRIGGER AS $$
BEGIN
   NEW.deleted_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_delete_bags
BEFORE DELETE ON bags
FOR EACH ROW
EXECUTE FUNCTION soft_delete_bags();