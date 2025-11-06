CREATE TABLE videos (
  video_id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  video_caption TEXT,
  upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  duration INT,
  video_url VARCHAR(512),
  thumbnail_url VARCHAR(512),
  likes_count INT DEFAULT 0,
  comments_count INT DEFAULT 0,
  views_count INT DEFAULT 0,
  is_public BOOLEAN DEFAULT TRUE,
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);
