CREATE TABLE users(
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    avatar_filename VARCHAR(100),
    avatar_url VARCHAR(255)
);