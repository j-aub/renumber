#!/bin/zsh
set -eu
# $0 is the server

# set timeformat to something parsable
TIMEFMT='%mE'
# ensure all files are cached
cat data/tiny > /dev/null
cat data/normal > /dev/null
cat data/massive > /dev/null
cat data/t-34.mp4 > /dev/null

echo 'testing tiny'
time curl -sS -F list=@data/tiny "$1" > /dev/null 2>&1

echo 'testing normal'
time curl -sS -F list=@data/normal "$1" > /dev/null 2>&1

echo 'testing massive'
time curl -sS -F list=@data/massive "$1" > /dev/null 2>&1

echo 'testing t-34.mp4'
time curl -sS -F list=@data/t-34.mp4 "$1" > /dev/null 2>&1

echo 'testing 1000 normal'
time (for i in $(seq 1 1000); do
curl -sS -F list=@data/normal "$1" > /dev/null 2>&1
done)
