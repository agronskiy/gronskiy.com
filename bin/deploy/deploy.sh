#!/bin/sh

set -e

(

    curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

    cd $curr_dir/../..

    if [ "`git status -s`" ]
    then
        echo "The working directory is dirty. Please commit any pending changes."
        exit 1;
    fi

    (
        cd content/posts/_unlisted
        if [ "`git status -s`" ]
        then
            echo "The '_unlisted' directory is dirty. Please commit any pending changes."
            exit 1;
        fi
    )

    # This pushes everything
    git push
    (
        cd content/posts/_unlisted
        git push
    )

    # This builds
    (
        echo "Checking out gh-pages branch into public"
        mkdir -p public
        cd public

        git fetch --all

        git reset --hard origin/gh-pages
        git checkout gh-pages
        git pull

        echo "Removing existing files"
        rm -rf *

        echo "Generating site"
        cd ..
        hugo

        echo "Updating gh-pages branch"
        cd public
        echo "gronskiy.com" > CNAME
        git add --all
        # Below commits only when something changed
        git diff-index --quiet HEAD | git commit -m "Publishing to gh-pages (publish.sh)"

        echo "Pushing public to github"
        git push
    )

)
