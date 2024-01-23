#!/bin/sh
# $1 takes an html template to split into header & footer
# it must have a single occurance of "{{ list }}" that we'll split on
set -eu

count="$(grep -c '{{ list }}' "$1")"
if [ "$count" -ne 1 ]; then
	echo 'wrong amount of "{{ list }}"'
	exit 1
fi

awk '/{{ list }}/{a++; sub(/{{ list }}.*/, ""); printf $0} !a{print $0}' "$1" > template/header.html
awk 'a{print $0} /{{ list }}/{a++; sub(/.*{{ list }}/, ""); print $0}' "$1" > template/footer.html
