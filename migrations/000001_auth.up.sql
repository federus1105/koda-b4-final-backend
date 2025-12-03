CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    role role DEFAULT 'user'
 );

CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    id_users INT NOT NULL,
    fullname VARCHAR(50),
    photos VARCHAR(100),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE account
ADD FOREIGN KEY (id_users) REFERENCES users(id);
CREATE INDEX idx_account_user_id ON account(user_id);