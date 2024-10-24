#!/bin/bash

: '
  Basic utility to upload csv file to postgres using `COPY`

    e.g.

      bash ./.scripts/upload-csv.sh -e .env -f ./.data/file.csv -d "|" -t "target_table(code, desc)"

'

usage()
{
  printf "Usage: $0\nOpts:\n\t[-f <string> content file path]\n\t[-d <string> field delimiter]\n\t[-t <string> target table name]\n\t[-e <string> env file containing pgsql connection info]\n" 1>&2;
  exit 1;
}

while getopts ":f:d:t:e:" flag
do
  case "${flag}" in
    f) file=${OPTARG};;
    d) delim=${OPTARG};;
    t) tabname=${OPTARG};;
    e) envfile=${OPTARG};;
    *)
      usage;;
  esac
done

# Assertions
if [ -z $file ]; then
  printf "No source file provided\n"
  usage
elif [ -z $tabname ]; then
  printf "No target table name provided\n"
  usage
fi

# Source env from file if provided
if [ ! -z envfile ]; then
  set -a && source "$envfile" && set +a
fi

# Defaults
if [ -z $delim ]; then
  delim="E'\\t'"
else
  export PGPASSWORD='$POSTGRES_PASSWORD'
fi

# Upload file
psql -d "$POSTGRES_DB" -U "$POSTGRES_USER" -p "$POSTGRES_PORT" -h "$POSTGRES_HOSTNAME" \
  -c "\copy $tabname from '$file' with (delimiter ${delim} format 'csv' encoding 'UTF-8' header 'on');"
