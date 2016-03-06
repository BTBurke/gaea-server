#!/bin/sh

postgres -D /var/lib/postgresql/data

psql -U postgres < /sql/db.sql


psql -U postgres db_gaea < /sql/schema.sql
