#!/bin/bash -e

env


if [ "$1" == "init" ]; then
    mkdir -p src/github.com/rdoorn
    cd src/github.com/rdoorn
    git clone https://github.com/rdoorn/ixxi.git
    cd ixxi
    make deps
    exit
fi

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
    bugfix:*|bug:*|fix:*)
        patch=$((patch+1))
        rebuild=1
        ;;
    feature:*|feat:*)
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

if [ $rebuild -eq 0 ]; then
    echo "version not updated: old: $(cat .version) new: ${major}.${minor}.${patch}"
    exit 0
fi

echo "new version to be created: old: $(cat .version) new: ${major}.${minor}.${patch}"
echo "${major}.${minor}.${patch}" > .version

sudo apt-get --no-install-recommends install ruby ruby-dev rubygems build-essential rpm
sudo gem install --no-ri --no-rdoc fpm

make linux-package

go get github.com/tcnksm/ghr
VERSION=$(cat .version)

ghr -soft -t ${GITHUB_TOKEN} -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} -c ${CIRCLE_SHA1} -n "${CIRCLE_PROJECT_REPONAME^} v${VERSION}" ${VERSION} ./build/packages/
