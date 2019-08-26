#!/bin/bash
env

major=$(cat .version | cut -f1 -d.)
minor=$(cat .version | cut -f2 -d.)
patch=$(cat .version | cut -f3 -d. | cut -f1 -d-)
case ${CIRCLE_BRANCH}
    bug-*)
        patch=$((patch+1))
        ;;
    feat-*)
        patch=0
        minor=$((minor+1))
        ;;
    major-*)
        patch=0
        minor=0
        major=$((major+1))
        ;;
esac

echo "old version: $(cat .version) new: ${major}.${minor}.${patch}"
echo "${major}.${minor}.${patch}" > .version
