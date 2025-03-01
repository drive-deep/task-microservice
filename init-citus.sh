#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$DATABASE_USER" --dbname "$DATABASE_NAME" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS citus;
    SELECT * FROM master_add_node('postgres-worker1', 5432);
    SELECT * FROM master_add_node('postgres-worker2', 5432);

    SELECT create_distributed_table('tasks', 'id');
EOSQL