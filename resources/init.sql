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
  nickname CITEXT COLLATE "ucs_basic" primary key
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

CREATE OR REPLACE FUNCTION ThreadInc() RETURNS trigger AS $$
BEGIN
  UPDATE forums SET threads = threads + 1
  WHERE slug = NEW.forum;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS ThreadInc ON Threads;

CREATE TRIGGER ThreadInc AFTER INSERT ON Threads
  FOR EACH ROW EXECUTE PROCEDURE ThreadInc();






CREATE TABLE IF NOT EXISTS votes (
  nickname  CITEXT     NOT NULL REFERENCES users(nickname),
  voice     INTEGER,
  thread    INTEGER     NOT NULL REFERENCES threads(id),
  UNIQUE(nickname, thread)
);

CREATE OR REPLACE FUNCTION VotesInc() RETURNS trigger AS $$
BEGIN
  UPDATE threads SET votes = votes + NEW.voice
  WHERE id = NEW.thread;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS VotesInc ON votes;

CREATE TRIGGER VotesInc AFTER INSERT ON votes
  FOR EACH ROW EXECUTE PROCEDURE VotesInc();

CREATE OR REPLACE FUNCTION VotesIncOnUpdate() RETURNS trigger AS $$
BEGIN
  UPDATE threads SET votes = votes + 2*NEW.voice
  WHERE id = NEW.thread;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS VotesIncOnUpdate ON votes;

CREATE TRIGGER VotesIncOnUpdate AFTER UPDATE ON votes
  FOR EACH ROW EXECUTE PROCEDURE VotesIncOnUpdate();






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


CREATE OR REPLACE FUNCTION postInc() RETURNS trigger AS $$
BEGIN
  UPDATE forums SET posts = posts + 1
  WHERE slug = NEW.forum;
  RETURN NEW;
END;
$$ language plpgsql;

DROP TRIGGER IF EXISTS postInc ON posts;

CREATE TRIGGER postInc AFTER insert ON posts
  FOR EACH ROW EXECUTE PROCEDURE postInc();

CREATE OR REPLACE FUNCTION addPath() RETURNS trigger
AS $$BEGIN
  NEW.path=array_append((SELECT path FROM posts WHERE id = NEW.parent), NEW.id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS addPathAfterInsertPost ON posts;

CREATE TRIGGER addPathAfterInsertPost
  BEFORE INSERT
  ON Posts
  FOR EACH ROW
EXECUTE PROCEDURE addPath();



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


DROP INDEX IF EXISTS boostUsernameIdx;
DROP INDEX IF EXISTS boostSlugIdx;
DROP INDEX IF EXISTS ThreadsIDIdx;
DROP INDEX IF EXISTS postsPathIdx;

CREATE INDEX IF NOT EXISTS threadsForumCreatedIdx ON Threads (forum, created);
CREATE INDEX IF NOT EXISTS votesUsernameThreadIdx ON Votes (nickname, thread);

CREATE INDEX IF NOT EXISTS postsThreadIdIdx ON Posts(thread, id, created);
CREATE INDEX IF NOT EXISTS postsPathIdx ON Posts(id,(path[1]));


CREATE TABLE IF NOT EXISTS Boost (
  username CITEXT NOT NULL REFERENCES Users (nickname),
  slug     CITEXT NOT NULL REFERENCES Forums (slug),
  UNIQUE (username, slug)
);



CREATE OR REPLACE FUNCTION addUserToBoost()
  RETURNS TRIGGER
AS $$
BEGIN
  INSERT INTO Boost (username, slug) VALUES (NEW.author, NEW.forum)
  ON CONFLICT DO NOTHING;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS AddUserToBoostAfterInsertOnThread ON Threads;

CREATE TRIGGER AddUserToBoostAfterInsertOnThread AFTER INSERT
  ON Threads
  FOR EACH ROW EXECUTE PROCEDURE addUserToBoost();

DROP TRIGGER IF EXISTS AddUserToBoostAfterInsertOnThread ON Posts;

CREATE TRIGGER AddUserToBoostAfterInsertOnPosts AFTER INSERT ON Posts
  FOR EACH ROW EXECUTE PROCEDURE addUserToBoost();

CREATE INDEX IF NOT EXISTS boostUsernameIdx ON Boost (username);
CREATE INDEX IF NOT EXISTS boostSlugIdx ON Boost (username, slug);
