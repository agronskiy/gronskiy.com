#!/bin/sh

set -e

(

    curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
    cd $curr_dir

    if [ "`git status -s`" ]
        then
            echo "The working directory is dirty. Please commit any pending changes."
            exit 1;
    fi

    echo "Generating tags"
    (
        ./generate-tags.sh
    )

    if [ "`git status -s`" ]
        then
            echo "Tags generation step changed working dir. Please commit any pending changes."
            exit 1;
    fi

    echo "Compiling sass"
    (
        ./compile-sass.sh
    )

    if [ "`git status -s`" ]
        then
            echo "SASS compilation step changed working dir. Please commit any pending changes."
            exit 1;
    fi

    echo "Building and deployment"
    (
        cd $curr_dir/deploy
        ./deploy.sh
    )

)
