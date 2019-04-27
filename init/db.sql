SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+03:00";

/* AUTH DB */

CREATE DATABASE IF NOT EXISTS colors_auth DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_auth;

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

CREATE DATABASE IF NOT EXISTS colors_core DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_core;

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

/* CHAT DB */

CREATE DATABASE IF NOT EXISTS colors_chat DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE colors_chat;

CREATE TABLE IF NOT EXISTS messages
(
    id      int(11)   NOT NULL AUTO_INCREMENT,
    created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_id int(11)   NOT NULL,
    type_id int(11)   NOT NULL,
    content text      NOT NULL,
    PRIMARY KEY (id),
    KEY type_id (type_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS types
(
    id          int(11)     NOT NULL AUTO_INCREMENT,
    type        varchar(20) NOT NULL,
    description text,
    PRIMARY KEY (id),
    UNIQUE KEY type (type)
) ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = utf8mb4;

ALTER TABLE messages
    ADD CONSTRAINT messages_ibfk_1 FOREIGN KEY (type_id) REFERENCES types (id);

COMMIT;

/* DATA */

# Passwords are same as usernames
INSERT INTO colors_auth.users (id, username, email, pass_hash)
VALUES (1, 'martini', 'martini@mail.ru', 'c0d542fcef483aaede5ec4c885205b30969f9545fab63e3a37d1304f7270eff2'),
       (2, 'jimbin', 'jimbin@mail.ru', 'edc50d35ca8a16b64eea128e69d33d7838f6b2916f72a93c3c58cb40fc7d687c'),
       (3, 'redlabel', 'redlabel@mail.ru', '14ef236e2db9f0f5ddf5d5cb248956cd7caa972ea1682f879a48a80b50c8ea3c'),
       (4, 'whitehorse', 'whitehorse@mail.ru', '05d62297e4f36c992e6e79b605d710ebdd3e3d7fe8a152f1d3309538440d5dbe'),
       (5, 'hennessy', 'hennessy@mail.ru', 'afa5731e35d576aa56a686ca3970a7b4e9ad5127dcbfdde35117949b9b3f53a6'),
       (6, 'ararat', 'ararat@mail.ru', 'd13039538bff2405191567278351171f1629a2127c157c86d35a84d21e347e5f'),
       (7, 'malibu', 'malibu@mail.ru', '5f16c3e2f49a617db5d9b49a3b433af11b054f8ee5700e3cb5b004c4cce20dca'),
       (8, 'jack_daniels', 'jack_daniels@mail.ru', 'e233bc06c2db20d5beab9c3945e78c3765ca519cda610bd053a84e8109345874'),
       (9, 'olmeca', 'olmeca@mail.ru', 'c7ca6655f679bcec5c6bec401813101813f18548af87fd2d6aa4a72b73e5ae40');

INSERT INTO colors_core.users (id, avatar, score, games, win, lose)
VALUES (1, '', 400, 10, 4, 6),
       (2, '', 600, 10, 6, 4),
       (3, '', 300, 10, 3, 7),
       (4, '', 200, 10, 2, 8),
       (5, '', 700, 10, 7, 3),
       (6, '', 100, 10, 1, 9),
       (7, '', 500, 10, 5, 5),
       (8, '', 800, 10, 8, 2),
       (9, '', 900, 10, 9, 1);

INSERT INTO colors_chat.types (id, type, description)
VALUES (1, 'text', 'Text message.'),
       (2, 'sticker', 'Sticker message.');

COMMIT;