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

echo "Deleting old publication"
rm -rf public
git worktree prune
rm -rf .git/worktrees/public/

echo "Checking out gh-pages branch into public"
mkdir public
git worktree add -B gh-pages public origin/gh-pages

cd public

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

#echo "Pushing to github"
git push --all

)
