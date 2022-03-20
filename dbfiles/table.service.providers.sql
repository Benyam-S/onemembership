CREATE TABLE service_providers (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    status INTEGER,
    created_at DATETIME,
    updated_at DATETIME
);