CREATE TABLE refresh_tokens (
  id SERIAL NOT NULL,
  refresh_token VARCHAR NOT NULL,
  user_id INT NOT NULL,
  user_agent VARCHAR NOT NULL,
  
  PRIMARY KEY (id),
  CONSTRAINT unique_refresh_token UNIQUE (refresh_token),
  CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);