CREATE TABLE "song_history"(
    "id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "song_id" BIGINT NOT NULL,
    "timestamp" DATE NOT NULL
);
ALTER TABLE
    "song_history" ADD PRIMARY KEY("id");
CREATE TABLE "songs_currently_playing"(
    "id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "song_id" BIGINT NOT NULL
);
ALTER TABLE
    "songs_currently_playing" ADD PRIMARY KEY("id");
CREATE TABLE "users"(
    "id" BIGINT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "avatar" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "available_songs"(
    "id" BIGINT NOT NULL,
    "track_name" VARCHAR(255) NOT NULL,
    "artists_name" VARCHAR(255) NOT NULL,
    "released_year" BIGINT NOT NULL
);
ALTER TABLE
    "available_songs" ADD PRIMARY KEY("id");
ALTER TABLE
    "available_songs" ADD CONSTRAINT "available_songs_track_name_unique" UNIQUE("track_name");
ALTER TABLE
    "songs_currently_playing" ADD CONSTRAINT "songs_currently_playing_song_id_foreign" FOREIGN KEY("song_id") REFERENCES "available_songs"("id");
ALTER TABLE
    "song_history" ADD CONSTRAINT "song_history_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "song_history" ADD CONSTRAINT "song_history_song_id_foreign" FOREIGN KEY("song_id") REFERENCES "available_songs"("id");
ALTER TABLE
    "songs_currently_playing" ADD CONSTRAINT "songs_currently_playing_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");