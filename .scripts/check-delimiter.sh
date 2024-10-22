#/bin/bash/

: '
  Basic utility to test tab/space delimination of an input file
'

declare -a delims=("\t" "," "|" "[[:space:]]")

usage() {
  printf "Usage: $0\nOpts:\n\t[-d <string> test all files in a directory]\n\t[-f <string> test specific file]\n" 1>&2;
  exit 0;
}

while getopts ":d:f:" flag
do
  case "${flag}" in
    d) dir=${OPTARG};;
    f) file=${OPTARG};;
    *)
      usage;;
  esac
done

space_or_tabs()
{
  content=$(head $1 -n 2)
  delimiter="unknown delimiter"

  for delim in "${delims[@]}"
  do
    test=$(echo "$content" | grep -q "$delim" && echo "1" || echo "0")
    if [ $test -eq 1 ]; then
      delimiter=$delim
      break
    fi
  done

  echo "Delimiter: ${delimiter}"
}

if [ -z "$dir" ] && [ -z "$file" ]; then
  usage
elif [ ! -z "$dir" ]; then
  if [ $(find "$dir" -type f | wc -l) -eq 0 ] ; then
    echo "No files to test in Dir<${dir}>"
    exit 0
  fi

  for file in $(find "$dir" -maxdepth 1 ! -type d); do
    space_or_tabs "$file"
  done
elif [ ! -z "$file" ]; then
  space_or_tabs "$file"
fi
