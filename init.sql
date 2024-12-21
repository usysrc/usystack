-- init.sql

-- Create the items table
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR NOT NULL
);

-- Seed the items if table is empty
INSERT INTO items (name) SELECT 'item1' WHERE NOT EXISTS (SELECT 1 FROM items WHERE name = 'item1');