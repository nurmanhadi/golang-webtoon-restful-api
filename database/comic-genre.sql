CREATE TABLE comic_genre(
    id int PRIMARY KEY AUTO_INCREMENT,
    comic_id VARCHAR(36) NOT NULL,
    genre_id INT NOT NULL,
    CONSTRAINT fk_comics_comic_genre FOREIGN KEY(comic_id) REFERENCES comics(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_genre_comic_genre FOREIGN KEY(genre_id) REFERENCES genres(id) ON DELETE CASCADE ON UPDATE CASCADE
);