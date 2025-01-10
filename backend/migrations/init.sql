-- Create the bags table with an additional "linked" column and "child_count" column
CREATE TABLE bags (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    bag_type VARCHAR(50) NOT NULL, -- "Parent" or "Child"
    child_count INT DEFAULT 0, -- Number of child bags for parent bags
    linked BOOLEAN DEFAULT FALSE, -- Indicates if the parent bag is linked to a bill
    parent_bag VARCHAR(255),
    deleted_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- -- Create the bag_map table for mapping parent and child bags
-- CREATE TABLE bag_map (
--     id SERIAL PRIMARY KEY,
--     parent_bag VARCHAR(255) NOT NULL,
--     child_bag VARCHAR(255) NOT NULL,
--     UNIQUE (parent_bag, child_bag),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- Create the links table for linking parent bags to bill IDs
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    parent_bag VARCHAR(255) NOT NULL,
    bill_id VARCHAR(255) NOT NULL,
    UNIQUE (parent_bag, bill_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for fast lookups
CREATE INDEX idx_bags_qr_code ON bags (qr_code);
CREATE INDEX idx_bags_linked ON bags (linked); -- Index for the "linked" field
-- CREATE INDEX idx_bag_map_parent_bag ON bag_map (parent_bag);
-- CREATE INDEX idx_bag_map_child_bag ON bag_map (child_bag);
CREATE INDEX idx_links_bill_id ON links (bill_id);
CREATE INDEX idx_parent_bag ON bags (parent_bag);

-- Soft delete function for the bags table
CREATE OR REPLACE FUNCTION soft_delete_bags()
RETURNS TRIGGER AS $$
BEGIN
   NEW.deleted_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for soft-deleting rows in the bags table
CREATE TRIGGER before_delete_bags
BEFORE DELETE ON bags
FOR EACH ROW
EXECUTE FUNCTION soft_delete_bags();