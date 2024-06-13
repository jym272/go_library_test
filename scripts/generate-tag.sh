#!/usr/bin/env bash

set -eou pipefail

# Set Git Configuration (Replace with your bot's info)
#git config --global user.email "your-bot-email@example.com"
#git config --global user.name "Your Bot Name"

last_tag=$(git describe --tags --abbrev=0)

IFS='.' read -r major minor patch <<< "$(echo "$last_tag" | tr -d 'v')"

echo -e "\033[32mLast tag: $last_tag\033[0m"

new_tag="v$major.$minor.$((patch+1))"

echo -e "\033[32mNew tag: $new_tag\033[0m"

git tag -a "$new_tag" -m "Release $new_tag"

git push origin "$new_tag"

echo -e "\033[35m✔ Done\033[0m"