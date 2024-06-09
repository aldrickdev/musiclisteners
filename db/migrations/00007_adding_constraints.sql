-- +goose Up
-- +goose StatementBegin
ALTER TABLE production.songs_currently_playing ADD CONSTRAINT "songs_currently_playing_song_id_foreign" FOREIGN KEY("song_id") REFERENCES production.available_songs("id");
ALTER TABLE production.songs_currently_playing ADD CONSTRAINT "songs_currently_playing_user_id_foreign" FOREIGN KEY("user_id") REFERENCES production.users("id");
ALTER TABLE production.song_history ADD CONSTRAINT "song_history_user_id_foreign" FOREIGN KEY("user_id") REFERENCES production.users("id");
ALTER TABLE production.song_history ADD CONSTRAINT "song_history_song_id_foreign" FOREIGN KEY("song_id") REFERENCES production.available_songs("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE production.songs_currently_playing DROP CONSTRAINT songs_currently_playing_song_id_foreign;
ALTER TABLE production.songs_currently_playing DROP CONSTRAINT songs_currently_playing_user_id_foreign;
ALTER TABLE production.song_history DROP CONSTRAINT song_history_user_id_foreign;
ALTER TABLE production.song_history DROP CONSTRAINT song_history_song_id_foreign;
-- +goose StatementEnd
