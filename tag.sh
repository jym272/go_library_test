#!/usr/bin/env bash

set -eou pipefail

# last tag in remote
last_tag=$(git describe --tags --abbrev=0)

# parse the tag i thre vars
IFS='.' read -r major minor patch <<< $(echo $last_tag | tr -d 'v')

# echo
echo "Last tag: $last_tag"
echo "Major: $major"
echo "Minor: $minor"
echo "Patch: $patch"

# create patch new tag
new_tag="v$major.$minor.$((patch+1))"
echo "New tag: $new_tag"

# Create the new tag
git tag -a $new_tag -m "Release $new_tag"

# Push the new tag

git push origin $new_tag