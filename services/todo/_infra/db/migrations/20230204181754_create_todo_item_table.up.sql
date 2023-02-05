CREATE TABLE IF NOT EXISTS todo_item(
    summary TEXT,
    id varchar(40) NOT NULL DEFAULT (uuid()) PRIMARY KEY, -- This can be improved by turning it into a binary
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted BOOL DEFAULT FALSE,
    completed BOOL DEFAULT FALSE
);