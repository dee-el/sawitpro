/**
  * This is the SQL script that will be used to initialize the database schema.
  * We will evaluate you based on how well you design your database.
  * 1. How you design the tables.
  * 2. How you choose the data types and keys.
  * 3. How you name the fields.
  * 
  * In this assignment we will use PostgreSQL as the database.
  */

CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);


CREATE TABLE user_attendance_summaries (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    total_login INT NOT NULL
);

CREATE TABLE user_attendance_logs (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    login_at TIMESTAMP NOT NULL
);
