#!/usr/bin/env bash

set -e

(
    curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

    cd $curr_dir/deploy
    ./deploy.sh
)
