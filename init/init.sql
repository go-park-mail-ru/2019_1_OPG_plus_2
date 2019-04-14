SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

CREATE DATABASE IF NOT EXISTS colors_auth DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_auth;

CREATE TABLE users (
  id int(11) NOT NULL,
  username varchar(32) NOT NULL,
  email varchar(128) NOT NULL,
  pass_hash varchar(64) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE users
  ADD PRIMARY KEY (id),
  ADD UNIQUE KEY username (username),
  ADD UNIQUE KEY email (email),
  ADD KEY username_2 (username,pass_hash),
  ADD KEY email_2 (email,pass_hash);

ALTER TABLE users
  MODIFY id int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;
COMMIT;

CREATE DATABASE IF NOT EXISTS colors_core DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_core;

CREATE TABLE users (
  id int(11) NOT NULL,
  avatar varchar(64) NOT NULL DEFAULT '',
  score int(11) NOT NULL DEFAULT '0',
  games int(11) NOT NULL DEFAULT '0',
  win int(11) NOT NULL DEFAULT '0',
  lose int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


ALTER TABLE users
  ADD PRIMARY KEY (id);
COMMIT;


/* DATA */

INSERT INTO colors_auth.users (id, username, email, pass_hash) VALUES
(1, 'test', 'test@mail.ru', '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08'), #password test
(2, 'user', 'user@mail.ru', '04f8996da763b7a969b1028ee3007569eaf3a635486ddab211d512c85b9df8fb'), #password user
(3, 'user1', 'user1@mail.ru', '0b14d501a594442a01c6859541bcb3e8164d183d32937b851835442f69d5c94e'), #password password1
(4, 'user2', 'user2@mail.ru', '6cf615d5bcaac778352a8f1f3360d23f02f34ec182e259897fd6ce485d7870d4'), #password password2
(5, 'user3', 'user3@mail.ru', '5906ac361a137e2d286465cd6588ebb5ac3f5ae955001100bc41577c3d751764'); #password password3

INSERT INTO colors_core.users (id, avatar, score, games, win, lose) VALUES
(1, '', 100500, 10, 7, 3),
(2, '', 1337, 12, 7, 5),
(3, '', 1982, 100, 20, 80),
(4, '', 100, 17, 7, 10),
(5, '', 228, 5, 2, 3);

COMMIT;