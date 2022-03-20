CREATE TABLE languages (
    code VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    flag BLOB,
    display_order INTEGER
);