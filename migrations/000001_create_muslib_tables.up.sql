CREATE TABLE IF NOT EXISTS artists(
   id SERIAL,
   name VARCHAR (64) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS songs(
   id SERIAL,
   artist_id INTEGER NOT NULL,
   title VARCHAR (64) NOT NULL,
   lyrics VARCHAR (8192),
   release_date VARCHAR (10),
   link VARCHAR (64),
   UNIQUE (artist_id, title)
);

CREATE INDEX artists_name_index ON artists (name);
CREATE INDEX songs_artist_id_index ON songs (artist_id);
CREATE INDEX songs_title_index ON songs (title);