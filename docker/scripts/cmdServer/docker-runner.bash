#!/usr/bin/env bash

bin/darwin/cmdServer -cmd="grep -i push"

while read line
do

  jq "$line"

done < "${1:-/dev/stdin}"
