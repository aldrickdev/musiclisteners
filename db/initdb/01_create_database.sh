#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 -U postgres -d postgres <<-EOSQL
  CREATE DATABASE musiclisteners;
EOSQL
