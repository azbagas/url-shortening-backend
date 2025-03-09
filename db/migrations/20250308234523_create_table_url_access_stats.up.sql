CREATE TABLE url_access_stats (
  id SERIAL NOT NULL,
  url_id INT NOT NULL,
  accessed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (id),
  CONSTRAINT fk_url_access_stats_url FOREIGN KEY (url_id) REFERENCES urls (id) ON DELETE CASCADE
);