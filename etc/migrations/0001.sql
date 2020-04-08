CREATE SCHEMA if NOT EXISTS test;

CREATE sequence if NOT EXISTS test.seq_users;
CREATE sequence if NOT EXISTS test.seq_comments;

CREATE TABLE if NOT EXISTS test.users
(
  id INT NOT NULL DEFAULT nextval('test.seq_users'::regclass),
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  CONSTRAINT "PK_users" PRIMARY KEY (id),
  CONSTRAINT "UQ_users_email" UNIQUE (email),
  CONSTRAINT "CHK_users_email" CHECK (email LIKE '%@%')
);

CREATE TABLE if NOT EXISTS test.comments
(
  id INT NOT NULL DEFAULT nextval('test.seq_comments'::regclass),
  id_user INT NOT NULL,
  txt VARCHAR NOT NULL,
  CONSTRAINT "PK_comments" PRIMARY KEY (id)
);