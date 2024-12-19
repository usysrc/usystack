-- init.sql

-- Create the items table
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR NOT NULL
);

-- Insert sample data
INSERT INTO items (name) VALUES ('Item 1');
INSERT INTO items (name) VALUES ('Item 2');
INSERT INTO items (name) VALUES ('Item 3');
