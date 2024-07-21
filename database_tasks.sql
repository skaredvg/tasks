CREATE DATABASE tasks;

CREATE TABLE users (
	id serial PRIMARY KEY,
	name text NOT NULL
);

CREATE TABLE tasks (
	id serial PRIMARY KEY,
	opened bigint NOT NULL DEFAULT extract(epoch from now()),
	closed bigint,
	author_id integer REFERENCES users(id) NOT NULL DEFAULT 0,
	assigned_id integer REFERENCES users(id) NOT NULL DEFAULT 0,
	title text NOT NULL,
	content text NOT NULL
);

CREATE TABLE labels (
	id serial PRIMARY KEY,
	name text NOT NULL
);

CREATE TABLE tasks_labels (
	tasks_id integer NOT NULL REFERENCES tasks(id),
	label_id integer NOT NULL REFERENCES labels(id)
)