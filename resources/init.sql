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

CREATE TABLE IF NOT EXISTS forums (
  posts INTEGER DEFAULT 0,
	slug CITEXT PRIMARY KEY,
	threads INTEGER DEFAULT 0,
	title CITEXT,
	author CITEXT REFERENCES users(nickname)
);


CREATE TABLE IF NOT EXISTS threads(
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


CREATE TABLE IF NOT EXISTS votes (
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


CREATE TABLE IF NOT EXISTS posts (
  author    CITEXT NOT NULL REFERENCES users(nickname),
  created   TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
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

DROP TRIGGER IF EXISTS add_path_after_insert_post ON posts;

CREATE TRIGGER add_path_after_insert_post
  BEFORE INSERT
  ON Posts
  FOR EACH ROW
EXECUTE PROCEDURE add_path();



DROP INDEX IF EXISTS usersLowerNicknameIdx;
DROP INDEX IF EXISTS usersLowerEmailIdx;
DROP INDEX IF EXISTS forumsNicknameIdx;
DROP INDEX IF EXISTS forumsSlugIdx;
DROP INDEX IF EXISTS threadsSlugIdx;
DROP INDEX IF EXISTS threadsAuthorIdx;
DROP INDEX IF EXISTS threadsForumIdx;
DROP INDEX IF EXISTS threadsForumCreatedIdx;
DROP INDEX IF EXISTS votesUsernameThreadIdx;
DROP INDEX IF EXISTS postsIdIdx;
DROP INDEX IF EXISTS postsAuthorIdx;
DROP INDEX IF EXISTS postsThreadIdx;
DROP INDEX IF EXISTS postsCreatedIdx;
DROP INDEX IF EXISTS postsForumIdx;
DROP INDEX IF EXISTS postsPath1Idx;
DROP INDEX IF EXISTS postsPath1ThreadIdx;
DROP INDEX IF EXISTS postsThreadIdIdx;



CREATE UNIQUE INDEX IF NOT EXISTS usersLowerNicknameIdx ON Users (LOWER(nickname));
CREATE UNIQUE INDEX IF NOT EXISTS usersLowerEmailIdx ON Users (LOWER(email));
CREATE INDEX IF NOT EXISTS forumsNicknameIdx ON Forums (LOWER(author));
CREATE UNIQUE INDEX IF NOT EXISTS forumsSlugIdx ON Forums (LOWER(slug));
CREATE INDEX IF NOT EXISTS threadsSlugIdx ON Threads (LOWER(slug));
CREATE INDEX IF NOT EXISTS threadsAuthorIdx on Threads (LOWER(author));
CREATE INDEX IF NOT EXISTS threadsForumIdx on Threads (LOWER(forum));
CREATE INDEX IF NOT EXISTS threadsForumCreatedIdx ON Threads (LOWER(forum), created);
CREATE INDEX IF NOT EXISTS votesUsernameThreadIdx ON Votes (nickname, thread);
CREATE INDEX IF NOT EXISTS postsIdIdx ON Posts (id);
CREATE INDEX IF NOT EXISTS postsAuthorIdx ON Posts (lower(author));
CREATE INDEX IF NOT EXISTS postsThreadIdx ON Posts (thread);
CREATE INDEX IF NOT EXISTS postsCreatedIdx ON Posts (created);
CREATE INDEX IF NOT EXISTS postsForumIdx ON Posts (lower(forum));
CREATE INDEX IF NOT EXISTS postsPath1Idx ON Posts ((path [1]));
CREATE INDEX IF NOT EXISTS postsPath1ThreadIdx ON Posts (thread, (path [1]));
CREATE INDEX IF NOT EXISTS postsThreadIdIdx ON Posts(thread, id);


CREATE TABLE IF NOT EXISTS Boost (
  username CITEXT NOT NULL REFERENCES Users (nickname),
  slug     CITEXT NOT NULL REFERENCES Forums (slug),
  UNIQUE (username, slug)
);



CREATE OR REPLACE FUNCTION addUserToBoost()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  INSERT INTO Boost (username, slug) VALUES (NEW.author, NEW.forum)
  ON CONFLICT DO NOTHING;
  RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS addUserToBoost ON posts;

CREATE TRIGGER add_user_after_insert_thread
  AFTER INSERT
  ON Threads
  FOR EACH ROW
EXECUTE PROCEDURE addUserToBoost();

CREATE TRIGGER add_user_after_insert_thread
  AFTER INSERT
  ON Posts
  FOR EACH ROW
EXECUTE PROCEDURE addUserToBoost();

CREATE INDEX IF NOT EXISTS boostUsernameIdx ON Boost (LOWER(username));
CREATE INDEX IF NOT EXISTS boostSlugIdx ON Boost (LOWER(slug), LOWER(username));
DROP INDEX IF EXISTS boostUsernameIdx;
DROP INDEX IF EXISTS boostSlugIdx;