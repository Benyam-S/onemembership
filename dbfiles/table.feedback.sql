CREATE TABLE feedbacks (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    comment BLOB NOT NULL,
    seen BOOLEAN,
    created_at DATETIME
);