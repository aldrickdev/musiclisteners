-- +goose Up
-- +goose StatementBegin
GRANT SELECT on production.users TO datadog;
GRANT SELECT on production.available_songs TO datadog;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
REVOKE SELECT on production.users TO datadog;
REVOKE SELECT on production.available_songs TO datadog;
-- +goose StatementEnd
