#!/bin/bash

env


echo "param1: $1"

#REF=$(git log --graph  --pretty=format:%D -1 | cut -f2 -d, | sed -e 's/.*\///g')
#echo "ref: ${REF}"
git describe --tags --always > .version
echo "path: ${PWD} version: $(cat .version)"

if [ "${CIRCLE_BRANCH}" != "master" ]; then
    echo "Branch is: [${CIRCLE_BRANCH}] skipping packaging"
    exit 0
fi

REF=$(git log -1 --pretty=%B)


major=$(cat .version | cut -f1 -d.)
minor=$(cat .version | cut -f2 -d.)
patch=$(cat .version | cut -f3 -d. | cut -f1 -d-)
rebuild=0
case "${REF}" in
    bugfix:*)
        patch=$((patch+1))
        rebuild=1
        ;;
    feature:*)
        patch=0
        minor=$((minor+1))
        rebuild=1
        ;;
    major:*)
        patch=0
        minor=0
        major=$((major+1))
        rebuild=1
        ;;
esac

if [ $rebuild -eq 1 ]; then
    echo "old version: $(cat .version) new: ${major}.${minor}.${patch}"
    echo "${major}.${minor}.${patch}" > .version


    gem install fpm
    make linux-package

    go get github.com/tcnksm/ghr
    VERSION=$(cat .version)
    echo ghr -t ${GITHUB_TOKEN} -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} -c ${CIRCLE_SHA1} -delete ${VERSION} ./build/packages/
fi
