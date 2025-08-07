CREATE TABLE comics(
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    synopsis TEXT NOT NULL,
    author VARCHAR(50) NOT NULL,
    artist VARCHAR(50) NOT NULL,
    type ENUM('manga', 'manhua', 'manhwa') NOT NULL,
    cover_filename VARCHAR(100),
    cover_url VARCHAR(255),
    views BIGINT DEFAULT NULL,
    updated_post TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);