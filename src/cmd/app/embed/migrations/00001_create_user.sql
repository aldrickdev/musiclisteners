-- +goose Up
-- +goose StatementBegin
-- +goose ENVSUB ON
CREATE USER app WITH PASSWORD '${APP_USER_POSTGRES_PASSWORD}';
-- +goose ENVSUB OFF
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP USER app;
-- +goose StatementEnd
