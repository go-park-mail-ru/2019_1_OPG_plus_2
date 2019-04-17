SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+03:00";

/* AUTH DB */

CREATE DATABASE IF NOT EXISTS colors_auth_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_auth_test;

CREATE TABLE users
(
    id        int(11)      NOT NULL,
    username  varchar(32)  NOT NULL,
    email     varchar(128) NOT NULL,
    pass_hash varchar(64)  NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

ALTER TABLE users
    ADD PRIMARY KEY (id),
    ADD UNIQUE KEY username (username),
    ADD UNIQUE KEY email (email),
    ADD KEY username_2 (username, pass_hash),
    ADD KEY email_2 (email, pass_hash);

ALTER TABLE users
    MODIFY id int(11) NOT NULL AUTO_INCREMENT,
    AUTO_INCREMENT = 3;

COMMIT;

/* CORE DB */

CREATE DATABASE IF NOT EXISTS colors_core_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_core_test;

CREATE TABLE users
(
    id     int(11)     NOT NULL,
    avatar varchar(64) NOT NULL DEFAULT '',
    score  int(11)     NOT NULL DEFAULT '0',
    games  int(11)     NOT NULL DEFAULT '0',
    win    int(11)     NOT NULL DEFAULT '0',
    lose   int(11)     NOT NULL DEFAULT '0'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

ALTER TABLE users
    ADD PRIMARY KEY (id);

COMMIT;
