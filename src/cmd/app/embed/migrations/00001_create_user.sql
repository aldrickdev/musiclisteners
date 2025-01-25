-- +goose Up
-- +goose StatementBegin
-- +goose ENVSUB ON
CREATE USER ${APP_USER_POSTGRES_USERNAME} WITH PASSWORD '${APP_USER_POSTGRES_PASSWORD}';
-- +goose ENVSUB OFF
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose ENVSUB ON
DROP USER ${APP_USER_POSTGRES_USERNAME} ;
-- +goose ENVSUB OFF
-- +goose StatementEnd
