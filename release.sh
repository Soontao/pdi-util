#!/bin/bash

if [ -z "$1" ]
then
  echo "must give out a new version"
else
  echo "this script will not commit current workspace un-commit files"
  git tag $1 >/dev/null
  git-chglog -o CHANGELOG.md >/dev/null
  git add -A >/dev/null
  git commit -m "chore(release): new version" >/dev/null
  git push --tags >/dev/null
  echo "new version $1 released"
fi

