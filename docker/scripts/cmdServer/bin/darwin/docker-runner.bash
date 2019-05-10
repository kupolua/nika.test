#!/usr/bin/env bash

while read line
do
    isGenerate=`echo "$line" | grep 'site built at'`

    if [[ -z "$isGenerate" ]]
     then
        echo "Nothing to do"
     else
         #git_url=`echo "$line" | sed -e 's/.*git_url\":\"\(.*\)\"/\1/p'`
         GIT_URL=`echo "$line" | jq '.repository.git_url'`
         GIT_EMAIL=`echo "$line" | jq '.head_commit.committer.email'`
         GIT_NAME=`echo "$line" | jq '.head_commit.committer.name'`

        docker run -it --rm -v ${HOME}/.ssh:/root/.ssh -e GIT_URL=${GIT_URL} -e GIT_EMAIL=${GIT_EMAIL} -e GIT_NAME=${GIT_NAME} kupolua/site-builder
    fi

done < "${1:-/dev/stdin}"


