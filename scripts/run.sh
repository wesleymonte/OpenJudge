#!/bin/bash

SCRIPT_PATH=$1
TESTS_PATH=$2

if [ -z "$SCRIPT_PATH" ] ; then
  echo "Script path may be not empty"
  exit 1
fi

if [ -z "$TESTS_PATH" ] ; then
  echo "Tests path may be not empty"
  exit 1
fi

mkdir -p results

for file in $(find "${TESTS_PATH}"/in/*.in)
do
    echo "$file"
    output=$(basename "$file")
    output=${output%%.*}.out
    python2 "${SCRIPT_PATH}" < "$file" >| ./results/"$output"
    echo ./results/"$output"
done

for file in $(find ./results/*.out)
do
    echo "$file"
    output=$(basename "$file")
    dif=$(diff -Z "$file" "${TESTS_PATH}"/out/"$output")
    echo "DIFF: $dif" 
    if [ -n "$dif" ]; then
        echo "WA" >| result
        echo "WA"
        break
    fi
    echo "AC" >| result
    echo "AC"
done
