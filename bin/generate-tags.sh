#!/bin/sh

set -e

(

curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

rm -rf ../data/tags_all.yaml
rm -rf ../tags/*

cd $curr_dir/tags-gen

go run main.go

)
