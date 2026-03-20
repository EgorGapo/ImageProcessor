-- Create tables
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(255) PRIMARY KEY,
    image_base TEXT NOT NULL,
    filter_name VARCHAR(255),
    filter_parametes TEXT,
    status VARCHAR(50),
    result TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS crawled_pages (
    url TEXT PRIMARY KEY
);
