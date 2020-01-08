DROP TABLE IF EXISTS Users CASCADE ;
DROP TABLE IF EXISTS forums CASCADE ;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS votes CASCADE ;
DROP TABLE IF EXISTS posts CASCADE ;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS Users (
  about    TEXT,
  email    CITEXT COLLATE "ucs_basic" NOT NULL UNIQUE,
  fullname VARCHAR,
  nickname CITEXT COLLATE "ucs_basic" NOT NULL UNIQUE
);

CREATE TABLE forums (
  posts INTEGER DEFAULT 0,
	slug CITEXT PRIMARY KEY,
	threads INTEGER DEFAULT 0,
	title CITEXT,
	author CITEXT REFERENCES users(nickname)
);


CREATE TABLE threads(
  author CITEXT REFERENCES users(nickname),
  created TIMESTAMP WITH TIME ZONE,
  forum CITEXT REFERENCES forums(slug),
	id BIGSERIAL PRIMARY KEY,
	message TEXT NOT NULL,
	slug CITEXT UNIQUE,
	title TEXT NOT NULL,
	votes BIGINT DEFAULT 0
);

CREATE OR REPLACE FUNCTION thread_inc() RETURNS trigger AS $$
BEGIN
  UPDATE forums SET threads = threads + 1
  WHERE slug = NEW.forum;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS thread_inc ON Threads;

CREATE TRIGGER thread_inc AFTER INSERT ON Threads
  FOR EACH ROW EXECUTE PROCEDURE thread_inc();


CREATE TABLE votes (
  nickname  CITEXT     NOT NULL REFERENCES users(nickname),
  voice     INTEGER,
  thread    INTEGER     NOT NULL REFERENCES threads(id),
  UNIQUE(nickname, thread)
);

CREATE OR REPLACE FUNCTION votes_inc() RETURNS trigger AS $$
BEGIN
  UPDATE threads SET votes = votes + NEW.voice
  WHERE id = NEW.thread;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS votes_inc ON votes;

CREATE TRIGGER votes_inc AFTER INSERT ON votes
  FOR EACH ROW EXECUTE PROCEDURE votes_inc();



CREATE OR REPLACE FUNCTION votes_inc_on_update() RETURNS trigger AS $$
BEGIN
  UPDATE threads SET votes = votes + 2*NEW.voice
  WHERE id = NEW.thread;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS votes_inc_on_update ON votes;

CREATE TRIGGER votes_inc_on_update AFTER UPDATE ON votes
  FOR EACH ROW EXECUTE PROCEDURE votes_inc_on_update();


CREATE TABLE posts (
  author    CITEXT NOT NULL REFERENCES users(nickname),
  created   TIMESTAMP WITH TIME ZONE,
  forum     CITEXT REFERENCES forums(slug),
  id        BIGSERIAL NOT NULL PRIMARY KEY,
  isEdited  BOOLEAN	DEFAULT FALSE,
  message   TEXT NOT NULL,
  parent    INTEGER	DEFAULT 0,
  thread    INTEGER	NOT NULL REFERENCES threads(id),
  path      BIGINT ARRAY
);


CREATE OR REPLACE FUNCTION post_inc() RETURNS trigger AS $$
BEGIN
  UPDATE forums SET posts = posts + 1
  WHERE slug = NEW.forum;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS post_inc ON posts;

CREATE TRIGGER post_inc AFTER insert ON posts
  FOR EACH ROW EXECUTE PROCEDURE post_inc();

CREATE OR REPLACE FUNCTION add_path() RETURNS trigger
LANGUAGE plpgsql
AS $$BEGIN
  NEW.path=array_append((SELECT path FROM posts WHERE id = NEW.parent), NEW.id);
  RETURN NEW;
END$$;

CREATE TRIGGER add_path_after_insert_post
  BEFORE INSERT
  ON Posts
  FOR EACH ROW
EXECUTE PROCEDURE add_path();