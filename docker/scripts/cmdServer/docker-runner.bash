#!/usr/bin/env bash

while read line
do

    echo 'kpl'
    #isGenerate=`echo "$line" | grep 'site built at'`

    #echo ${isGenerate} 'kpl'

done < "${1:-/dev/stdin}"
