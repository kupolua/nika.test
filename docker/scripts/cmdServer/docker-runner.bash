#!/usr/bin/env bash

bin/darwin/cmdServer -cmd="grep -i push"

while read line
do

  echo "$line"

done < "${1:-/dev/stdin}"
