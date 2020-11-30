#!/bin/sh

set -e

(

curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $curr_dir

sass -t compressed ../scss/main.scss ../themes/gron/static/css/main.min.css
sass --scss ../scss/main.scss ../themes/gron/static/css/main.css

)


