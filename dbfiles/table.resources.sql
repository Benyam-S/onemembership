CREATE TABLE resources (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    name BLOB NOT NULL,
    description BLOB,
    type VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (provider_id) REFERENCES service_providers(id) ON DELETE CASCADE ON UPDATE CASCADE
);
