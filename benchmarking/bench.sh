#!/bin/zsh
set -eu
# $0 is the server

echo 'testing tiny'
time curl -F list=@data/tiny "$1" >/dev/null 2>&1

echo 'testing normal'
time curl -F list=@data/normal "$1" >/dev/null 2>&1

echo 'testing massive'
time curl -F list=@data/massive "$1" >/dev/null 2>&1

echo 'testing t-34.mp4'
time curl -F list=@data/t-34.mp4 "$1" >/dev/null 2>&1

echo 'testing 1000 normal'
time (for i in $(seq 1 1000); do
curl -F list=@data/normal "$1" >/dev/null 2>&1
done)
