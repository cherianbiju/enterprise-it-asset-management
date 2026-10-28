-- =============================================================================
-- IT Asset Management Database Schema
-- Database: PostgreSQL
-- =============================================================================

-- Drop tables if they already exist (order matters due to foreign keys)
DROP TABLE IF EXISTS asset_history CASCADE;
DROP TABLE IF EXISTS assets CASCADE;

-- -----------------------------------------------------------------------------
-- 1. Table: assets
-- -----------------------------------------------------------------------------
CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    asset_code VARCHAR(50) NOT NULL UNIQUE,
    asset_name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(50) NOT NULL,
    serial_number VARCHAR(100) NOT NULL UNIQUE,
    purchase_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (status IN ('AVAILABLE', 'ASSIGNED', 'UNDER_REPAIR', 'RETIRED')),
    assigned_to VARCHAR(100),
    location VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------------
-- 2. Table: asset_history
-- -----------------------------------------------------------------------------
CREATE TABLE asset_history (
    id SERIAL PRIMARY KEY,
    asset_id INT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    performed_by VARCHAR(100) DEFAULT 'System',
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------------
-- Indexes for Optimized Query Performance
-- -----------------------------------------------------------------------------
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_type ON assets(asset_type);
CREATE INDEX idx_assets_code ON assets(asset_code);
CREATE INDEX idx_asset_history_asset_id ON asset_history(asset_id);
CREATE INDEX idx_asset_history_created_at ON asset_history(created_at DESC);

-- -----------------------------------------------------------------------------
-- Sample Seed Data
-- -----------------------------------------------------------------------------
INSERT INTO assets (asset_code, asset_name, asset_type, serial_number, purchase_date, status, assigned_to, location)
VALUES 
    ('AST-1001', 'MacBook Pro 16" M3 Max', 'LAPTOP', 'C02G4567MD6R', '2024-01-15', 'ASSIGNED', 'Sarah Connor', 'San Francisco HQ - 4th Floor'),
    ('AST-1002', 'Dell XPS 15 9530', 'LAPTOP', 'DL9874561230', '2024-02-20', 'AVAILABLE', NULL, 'IT Storage Room B'),
    ('AST-1003', 'Dell UltraSharp 27" 4K', 'MONITOR', 'MON-DEL-4589', '2023-11-10', 'ASSIGNED', 'Sarah Connor', 'San Francisco HQ - 4th Floor'),
    ('AST-1004', 'Lenovo ThinkPad X1 Carbon', 'LAPTOP', 'PF29KLM89012', '2023-09-05', 'UNDER_REPAIR', NULL, 'Service Center - TechFix'),
    ('AST-1005', 'HP ProLiant DL380 Gen10', 'SERVER', 'SVR-HP-887711', '2022-06-18', 'ASSIGNED', 'DevOps Infrastructure', 'Data Center Rack 04'),
    ('AST-1006', 'iPhone 15 Pro 256GB', 'PHONE', 'APL-IPH15-9988', '2024-03-01', 'AVAILABLE', NULL, 'IT Secure Cabinet 1'),
    ('AST-1007', 'Apple Studio Display 27"', 'MONITOR', 'APL-STD-77665', '2023-08-12', 'AVAILABLE', NULL, 'New York Office - Pod C'),
    ('AST-1008', 'Cisco Catalyst 9300 Switch', 'NETWORK', 'CSCO-SW-33214', '2022-04-10', 'ASSIGNED', 'Network Operations', 'Data Center Core Switch'),
    ('AST-1009', 'MacBook Pro 14" M1 Pro', 'LAPTOP', 'C02F99887766', '2021-12-01', 'RETIRED', NULL, 'E-Waste Disposal Bin'),
    ('AST-1010', 'Logitech MX Master 3S', 'PERIPHERAL', 'LOGI-MX3S-1122', '2024-02-15', 'AVAILABLE', NULL, 'IT Accessories Shelf');

-- Seed History Logs for initial tracking
INSERT INTO asset_history (asset_id, action, performed_by, details, created_at)
VALUES 
    (1, 'CREATED', 'Admin', 'Asset initialized into inventory with status AVAILABLE', '2024-01-15 10:00:00+00'),
    (1, 'ASSIGNED', 'Admin', 'Assigned asset to Sarah Connor in San Francisco HQ - 4th Floor', '2024-01-16 14:30:00+00'),
    (2, 'CREATED', 'Admin', 'Asset initialized into inventory with status AVAILABLE', '2024-02-20 09:15:00+00'),
    (3, 'CREATED', 'Admin', 'Asset initialized into inventory with status AVAILABLE', '2023-11-10 11:00:00+00'),
    (3, 'ASSIGNED', 'Admin', 'Assigned monitor to Sarah Connor', '2023-11-12 16:00:00+00'),
    (4, 'CREATED', 'Admin', 'Asset initialized into inventory with status AVAILABLE', '2023-09-05 08:30:00+00'),
    (4, 'ASSIGNED', 'Admin', 'Assigned to Alex Morgan', '2023-09-06 10:00:00+00'),
    (4, 'RETURNED', 'Admin', 'Returned by Alex Morgan due to battery and thermal throttling issues', '2024-03-10 09:00:00+00'),
    (4, 'REPAIR', 'Admin', 'Transferred to Service Center - TechFix for diagnostics and replacement', '2024-03-11 11:20:00+00'),
    (9, 'CREATED', 'Admin', 'Asset initialized into inventory', '2021-12-01 10:00:00+00'),
    (9, 'ASSIGNED', 'Admin', 'Assigned to Engineering intern', '2021-12-05 14:00:00+00'),
    (9, 'RETURNED', 'Admin', 'Returned asset at conclusion of internship', '2023-12-20 15:30:00+00'),
    (9, 'RETIRED', 'Admin', 'Decommissioned: end of life, motherboard degradation. Cleared for disposal', '2024-01-05 16:45:00+00');
