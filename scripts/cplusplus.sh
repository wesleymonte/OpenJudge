#!/bin/bash

SCRIPT_PATH=$1
TESTS_PATH=$2

FILE_NAME_EXT=$(basename "$SCRIPT_PATH")
FILE_NAME=${FILE_NAME_EXT%.*}

mkdir -p results

for file in $(find "${TESTS_PATH}"/in/*.in)
do
    echo "$file"
    output=$(basename "$file")
    output=${output%%.*}.out
    g++ "$SCRIPT_PATH" -o "$FILE_NAME"
    ./"$FILE_NAME" < "$file" >| ./results/"$output"
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




