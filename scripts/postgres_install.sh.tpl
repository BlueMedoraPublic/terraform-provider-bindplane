#!/bin/bash


yum install postgresql-server postgresql-contrib -y


postgresql-setup initdb


cat << EOF > /var/lib/pgsql/data/pg_hba.conf
# Database administrative login by Unix domain socket
local   all             postgres                                peer
# TYPE  DATABASE        USER            ADDRESS                 METHOD
# "local" is for Unix domain socket connections only
local   all             all                                     peer
# IPv4 local connections:i
host    all             all             0.0.0.0/0               md5
# IPv6 local connections:
host    all             all             ::1/128                 md5
# Allow replication connections from localhost, by a user with the
# replication privilege.
#local   replication     postgres                                peer
#host    replication     postgres        127.0.0.1/32            md5
#host    replication     postgres        ::1/128                 md5
EOF

echo "shared_preload_libraries = 'pg_stat_statements'" >> /var/lib/pgsql/data/postgresql.conf
echo "track_activity_query_size = 10000" >> /var/lib/pgsql/data/postgresql.conf
echo "pg_stat_statements.track = all" >> /var/lib/pgsql/data/postgresql.conf
echo "listen_addresses = '*'" >> /var/lib/pgsql/data/postgresql.conf


systemctl enable postgresql
systemctl start postgresql


cat << EOF > /tmp/setup.sql
CREATE DATABASE app;
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
\c app
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
CREATE USER "${username}" WITH PASSWORD '${password}';
GRANT all privileges ON DATABASE app TO "${username}";
EOF

su - postgres -c "psql -f /tmp/setup.sql"


systemctl restart postgresql
