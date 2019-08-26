#!/bin/bash

go get github.com/tcnksm/ghr
VERSION=$(cat .version)
echo ghr -t ${GITHUB_TOKEN} -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} -c ${CIRCLE_SHA1} -delete ${VERSION} ./build/packages/

