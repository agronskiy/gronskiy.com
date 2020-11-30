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

        echo "Checking out gh-pages branch into public"
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
        git commit -m "Publishing to gh-pages (publish.sh)"

        echo "Pushing public to github"
        git push
    )

    # This pushes everything
    git push --recurse-submodules=on-demand

)
