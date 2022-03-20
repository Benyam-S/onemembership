CREATE TABLE deleted_users (
    id INTEGER PRIMARY KEY UNIQUE AUTO_INCREMENT NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    phone_number VARCHAR(255) NOT NULL,
    created_at DATETIME
);