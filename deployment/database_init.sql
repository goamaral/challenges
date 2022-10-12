CREATE TABLE users (
  id CHAR(26) NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  first_name VARCHAR(32) NOT NULL,
  last_name VARCHAR(32) NOT NULL,
  nickname VARCHAR(32) NOT NULL UNIQUE,
  encrypted_password CHAR(60) NOT NULL,
  email VARCHAR(64) NOT NULL UNIQUE,
	country VARCHAR(32) NOT NULL
);