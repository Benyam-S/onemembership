CREATE TABLE subscription_plans (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    owner_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description BLOB,
    duratinon INTEGER NOT NULL,
    price DOUBLE PRECISION(11, 2) NOT NULL,
    currency VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (owner_id) REFERENCES service_providers(id) ON DELETE CASCADE ON UPDATE CASCADE
);