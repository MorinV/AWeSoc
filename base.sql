CREATE DATABASE AWeSoc CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE AWeSoc;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    login VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    created DATETIME NOT NULL,
    last_login DATETIME NOT NULL
);

CREATE TABLE personal (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    firstname VARCHAR(100) NOT NULL,
    secondname VARCHAR(100),
    surname VARCHAR(100) NOT NULL,
    fullname VARCHAR(310) NOT NULL,
    birthdate DATETIME NOT NULL,
    gender ENUM('Мужчина', 'Женщина'),
    city VARCHAR(255) NOT NULL,
    interests text,
    user_id INTEGER
);

CREATE TABLE friends (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    personal_id INTEGER NOT NULL,
    friend_id INTEGER NOT NULL,
    state ENUM('pending', 'approved'),
    created DATETIME NOT NULL
);

CREATE UNIQUE INDEX UNQ_users_login ON users (login);
CREATE UNIQUE INDEX UNQ_friends ON friends (personal_id, friend_id);
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE ON AWeSoc.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'qwerty'