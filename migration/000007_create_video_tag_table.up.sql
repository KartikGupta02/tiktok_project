CREATE TABLE video_tags (
  video_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (video_id, tag_id),
  FOREIGN KEY (video_id) REFERENCES videos(video_id),
  FOREIGN KEY (tag_id) REFERENCES tags(tag_id)
);
