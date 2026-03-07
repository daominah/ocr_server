#!/bin/bash

# this script read loop through all files *.png in the current directory
# and store "data:image/png;base64ed" to a file name "*.base64.txt"

for f in *.png; do
    echo -n "data:image/png;base64," >"${f%.png}.base64.txt"
    base64 -w 0 "$f" >>"${f%.png}.base64.txt"
    echo "wrote ${f%.png}.base64.txt"
done
