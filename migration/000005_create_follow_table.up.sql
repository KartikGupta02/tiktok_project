CREATE TABLE follows (
  follow_id INT AUTO_INCREMENT PRIMARY KEY,
  follower_user_id INT NOT NULL,
  following_user_id INT NOT NULL,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (follower_user_id) REFERENCES users(user_id),
  FOREIGN KEY (following_user_id) REFERENCES users(user_id)
);
