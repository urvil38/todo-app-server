#!/usr/bin/env sh

GIT_SHA=$(git rev-parse HEAD)

if [ -z "$(git status --porcelain 2>/dev/null)" ]; then
    echo $GIT_SHA
else
    echo "$GIT_SHA-dirty"
fi
