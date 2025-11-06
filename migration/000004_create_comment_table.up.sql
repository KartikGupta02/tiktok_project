CREATE TABLE comments (
  comment_id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  video_id INT NOT NULL,
  comment_text TEXT,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  parent_comment_id INT DEFAULT NULL,
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (video_id) REFERENCES videos(video_id),
  FOREIGN KEY (parent_comment_id) REFERENCES comments(comment_id)
);