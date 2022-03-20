CREATE TABLE deleted_service_providers (
    id INTEGER PRIMARY KEY UNIQUE AUTO_INCREMENT NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    created_at DATETIME
);