CREATE TABLE language_entries (
    id INTEGER PRIMARY KEY UNIQUE AUTO_INCREMENT NOT NULL,
    identifier VARCHAR(255) NOT NULL,
    value BLOB NOT NULL,
    code VARCHAR(255) NOT NULL,
    UNIQUE KEY unique_language_entry (identifier, code),
    FOREIGN KEY (code) REFERENCES languages(code) ON DELETE CASCADE ON UPDATE CASCADE
);