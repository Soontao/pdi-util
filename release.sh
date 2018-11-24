#!/bin/bash

if [ -z "$1" ]
then
  echo "must give out a new version"
else
  git add -A
  git commit -m "chore(release): new version"
  git tag $1
  git-chglog -o CHANGELOG.md
  git add -A
  git commit -m "chore(release): update CHANGELOG"
  git push --tags
  echo "new version $1 released"
fi

