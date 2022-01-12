#!/usr/bin/env sh

usage() {
  cat <<EOUSAGE
Usage: $0 [up|down|force|version] {#}"
EOUSAGE
}

database_user="postgres"
if [[ $TODO_DATABASE_USER != "" ]]; then
  database_user=$TODO_DATABASE_USER
fi
database_password=""
if [[ $TODO_DATABASE_PASSWORD != "" ]]; then
  database_password=$TODO_DATABASE_PASSWORD
fi
database_host="localhost"
if [[ $TODO_DATABASE_HOST != "" ]]; then
  database_host=$TODO_DATABASE_HOST
fi
database_name='todo-db'
if [[ $TODO_DATABASE_NAME != "" ]]; then
  database_name=$TODO_DATABASE_NAME
fi

# Redirect stderr to stdout because migrate outputs to stderr, and we want
# to be able to use ordinary output redirection.
case "$1" in
  up|down|force|version)
    migrate \
      -source file:migrations \
      -database "postgresql://$database_user:$database_password@$database_host:5432/$database_name?sslmode=disable" \
      "$@" 2>&1
    ;;
  *)
    usage
    exit 1
    ;;
esac
