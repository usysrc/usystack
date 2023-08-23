-- init.sql

-- Create the items table
CREATE TABLE items (
    id serial PRIMARY KEY,
    name VARCHAR NOT NULL
);

-- Insert sample data
INSERT INTO items (name) VALUES ('Item 1');
INSERT INTO items (name) VALUES ('Item 2');
INSERT INTO items (name) VALUES ('Item 3');
