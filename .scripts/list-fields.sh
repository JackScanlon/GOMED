#!/bin/bash

: '
  Basic utility to find/list fields in a delimited file

    e.g.

      bash ./.scripts/list-fields.sh -f .data/some_file.txt -n 10 -c '$1 == "00589"' -d '[[:space:]]'

'

usage()
{
  printf "Usage: $0\nOpts:\n\t[-f <string> file containing rows]\n\t[-d <string> field delimiter]\n\t[-n <int|0> max number of matches before exit; where 0 = infinite matches]\n\t[-c <string> value comparator condition]\n" 1>&2;
  exit 1;
}

while getopts ":f:d:n:c:" flag
do
  case "${flag}" in
    f) file=${OPTARG};;
    d) delim=${OPTARG};;
    n) maxcount=$((${OPTARG} + 0));;
    c) comparator=${OPTARG};;
    *)
      usage;;
  esac
done

if [ -z $file ]; then
  printf "Missing valid filepath\n"
  usage
elif [ -z "$comparator" ]; then
  printf "Missing value comparator condition\n"
  usage
fi

if [ -z $delim ]; then
  delim="\t"
fi

if [ -z $maxcount ]; then
  maxcount=0
fi

awk -F"$delim" -v OFS="$delim" -v rcnt=0 -v mcnt=$maxcount '
  {
    if ('"${comparator}"') {
      rcnt++;
      print $0;
    }

    if (mcnt > 0 && rcnt >= mcnt) {
      exit
    }
  }' $file
