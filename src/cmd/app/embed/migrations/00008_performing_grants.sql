-- +goose Up
-- +goose StatementBegin
GRANT ALL PRIVILEGES ON DATABASE musiclisteners TO app;
GRANT USAGE on SCHEMA production TO app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA production TO app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA production TO app;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
REVOKE USAGE, SELECT ON ALL SEQUENCES IN SCHEMA production FROM app;
REVOKE SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA production FROM app;
REVOKE USAGE ON SCHEMA production FROM app;
REVOKE ALL PRIVILEGES ON DATABASE musiclisteners FROM app;
-- +goose StatementEnd
