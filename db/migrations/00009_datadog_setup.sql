-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA datadog;

GRANT USAGE ON SCHEMA datadog TO datadog;
GRANT USAGE ON SCHEMA public TO datadog;
GRANT USAGE ON SCHEMA production TO datadog;
GRANT pg_monitor TO datadog;
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE OR REPLACE FUNCTION datadog.explain_statement(
   l_query TEXT,
   OUT explain JSON
)
RETURNS SETOF JSON AS
$$
DECLARE
curs REFCURSOR;
plan JSON;

BEGIN
   OPEN curs FOR EXECUTE pg_catalog.concat('EXPLAIN (FORMAT JSON) ', l_query);
   FETCH curs INTO plan;
   CLOSE curs;
   RETURN QUERY SELECT plan;
END;
$$
LANGUAGE 'plpgsql'
RETURNS NULL ON NULL INPUT
SECURITY DEFINER;

ALTER ROLE datadog SET search_path = "$user",public,production;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS datadog.explain_statement(TEXT);

DROP EXTENSION IF EXISTS pg_stat_statements;
REVOKE pg_monitor FROM datadog;
REVOKE USAGE ON SCHEMA datadog FROM datadog;
REVOKE USAGE ON SCHEMA public FROM datadog;
REVOKE USAGE ON SCHEMA production TO datadog;
DROP SCHEMA IF EXISTS datadog CASCADE;

ALTER ROLE datadog NOINHERIT;

DROP USER IF EXISTS datadog;
-- +goose StatementEnd
