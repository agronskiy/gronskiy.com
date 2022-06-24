#!/usr/bin/env bash

set -e

# Check both branches are neither master nor main
(
    curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

    branch="$(git symbolic-ref --short HEAD)"
    if [ "$branch" == "main" ] \
        || [ "$branch" == "master" ];
    then
        echo "Branch is main/master!"
        exit 1
    fi

    cd $curr_dir/../content/posts/_unlisted

    branch="$(git symbolic-ref --short HEAD)"
    if [ "$branch" == "main" ] \
        || [ "$branch" == "master" ];
    then
        echo "Branch _unlisted is main/master!"
        exit 1
    fi

)

# Commit both branches
(
    curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

    if [ "`git status -s`" ]
    then
        git commit -am "WIP"
    fi

    cd $curr_dir/../content/posts/_unlisted

    if [ "`git status -s`" ]
    then
        git commit -am "WIP"
    fi
)


(
    curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

    cd $curr_dir/deploy
    ./deploy.sh
)
