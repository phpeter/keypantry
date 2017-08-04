/*
Users:
	-id: serial
	-username: character varying(55)
	-passwordHash: character varying(255)

Items:
	-id: serial
	-name: character varying(55)
	-userId: int
	-key: int
	-own: boolean
*/

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username CHARACTER VARYING(55) NOT NULL,
	passwordHash CHARACTER VARYING(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	name CHARACTER VARYING(55),
	userId INT REFERENCES users,
	key INT,
	isOwned BOOLEAN
);

create table IF NOT EXISTS usersession (
 SessionKey text primary key,
 UserID int not null REFERENCES users,
 LoginTime timestamp not null,
 LastSeenTime timestamp not null
)


INSERT INTO users (username, passwordHash) VALUES ('tindell', '6ca13d52ca70c883e0f0bb101e425a89e8624de51db2d2392593af6a84118090');
INSERT INTO items (name, userId, key) VALUES ('Cilantro', 1, 97);