-- Create bills table
CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    bill_code VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create bags table
CREATE TABLE IF NOT EXISTS bags (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    bag_type VARCHAR(50) NOT NULL CHECK (bag_type IN ('Parent', 'Child')),
    child_count INT DEFAULT 0,
    linked BOOLEAN DEFAULT FALSE,
    linked_to_bill BOOLEAN DEFAULT FALSE,
    parent_bag_id INT REFERENCES bags(id) ON DELETE SET NULL,
    bill_id INT REFERENCES bills(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create links table
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    parent_bag_id INT REFERENCES bags(id) ON DELETE CASCADE,
    bill_id INT REFERENCES bills(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(parent_bag_id, bill_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_bags_qr_code    ON bags(qr_code);
CREATE INDEX IF NOT EXISTS idx_bags_linked     ON bags(linked);
CREATE INDEX IF NOT EXISTS idx_bags_parent_id  ON bags(parent_bag_id);
CREATE INDEX IF NOT EXISTS idx_bags_bill_id    ON bags(bill_id);
CREATE INDEX IF NOT EXISTS idx_links_bill_id   ON links(bill_id);
CREATE INDEX IF NOT EXISTS idx_bills_bill_code ON bills(bill_code);

-- Trigger function to auto-update the updated_at timestamp on bags
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to set updated_at before each update on bags
DROP TRIGGER IF EXISTS trg_set_bags_updated ON bags;
CREATE TRIGGER trg_set_bags_updated
BEFORE UPDATE ON bags
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Trigger function to implement soft delete on bags (set deleted_at instead of delete)
CREATE OR REPLACE FUNCTION soft_delete_bag() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.deleted_at IS NULL THEN
        -- Mark the bag as deleted by setting deleted_at timestamp
        UPDATE bags SET deleted_at = NOW() WHERE id = OLD.id;
    END IF;
    -- Skip physical deletion of the row
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to invoke soft delete before a bag is deleted
DROP TRIGGER IF EXISTS trg_soft_delete_bag ON bags;
CREATE TRIGGER trg_soft_delete_bag
BEFORE DELETE ON bags
FOR EACH ROW
EXECUTE FUNCTION soft_delete_bag();