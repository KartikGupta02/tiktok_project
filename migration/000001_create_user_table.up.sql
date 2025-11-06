CREATE TABLE users (
  user_id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  profile_picture_url VARCHAR(512),
  bio VARCHAR(1024),
  followers_count INT DEFAULT 0,
  followings_count INT DEFAULT 0,
  registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
ADD CONSTRAINT unique_username UNIQUE (username);

