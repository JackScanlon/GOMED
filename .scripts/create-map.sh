#!/bin/bash

: '
  Basic utility to export columns from mapping data file

    e.g.

      bash ./.scripts/create-map.sh -i .data/some_file.txt -o .data/some_output.txt -c '$1$2$3' -d "|"

'

usage()
{
  printf "Usage: $0\nOpts:\n\t[-i <string> input file]\n\t[-o <string> output file]\n\t[-d <string> field delimiter]\n\t[-c <string> columns e.g. "'$1,$2,$3'"]\n" 1>&2;
  exit 1;
}

while getopts ":i:o:d:c:" flag
do
  case "${flag}" in
    i) input=${OPTARG};;
    o) output=${OPTARG};;
    d) delim=${OPTARG};;
    c) columns=${OPTARG};;
    *)
      usage;;
  esac
done

if [ -z $input ]; then
  printf "Missing valid input filepath\n"
  usage
elif [ -z $output ]; then
  printf "Missing valid output filepath\n"
  usage
elif [ -z "$columns" ]; then
  printf "Missing column selector\n"
  usage
fi

if [ -z $delim ]; then
  delim="\t"
fi

awk -F"$delim" -v OFS="$delim" -v rcnt=0 -v mcnt=$maxcount '
  {
    print '"${columns}"' > "'"${output}"'"
  }' $input
