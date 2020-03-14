# This file is part of the Smart Home
# Program complex distribution https://github.com/e154/smart-home
# Copyright (C) 2016-2020, Filippov Alex
#
# This library is free software: you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public
# License as published by the Free Software Foundation; either
# version 3 of the License, or (at your option) any later version.
#
# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
# Library General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public
# License along with this library.  If not, see
# <https://www.gnu.org/licenses/>.

#!/usr/bin/env bash

set -o errexit

#
# base variables
#
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
EXEC="server"
TMP_DIR="${ROOT}/tmp/${EXEC}"
ARCHIVE="smart-home-${EXEC}.tar.gz"

VERSION_VALUE="$(git describe --always --dirty --tags 2>/dev/null)"


#
# docker params
#
DEPLOY_IMAGE=smart-home-${EXEC}
DOCKER_VERSION="${VERSION_VALUE//-dirty}"
IMAGE=smart-home-${EXEC}
DOCKER_ACCOUNT=e154
DOCKER_IMAGE_VER=${DOCKER_ACCOUNT}/${IMAGE}:${DOCKER_VERSION}
DOCKER_IMAGE_LATEST=${DOCKER_ACCOUNT}/${IMAGE}:latest

#
# build version variables
#
PACKAGE="github.com/e154/smart-home"
VERSION_VAR="${PACKAGE}/version.VersionString"
REV_VAR="${PACKAGE}/version.RevisionString"
REV_URL_VAR="${PACKAGE}/version.RevisionURLString"
GENERATED_VAR="${PACKAGE}/version.GeneratedString"
DEVELOPERS_VAR="${PACKAGE}/version.DevelopersString"
BUILD_NUMBER_VAR="${PACKAGE}/version.BuildNumString"
DOCKER_IMAGE_VAR="${PACKAGE}/version.DockerImageString"
REV_VALUE="$(git rev-parse HEAD 2> /dev/null || echo "???")"
REV_URL_VALUE="https://${PACKAGE}/commit/${REV_VALUE}"
GENERATED_VALUE="$(date -u +'%Y-%m-%dT%H:%M:%S%z')"
DEVELOPERS_VALUE="delta54<support@e154.ru>"
BUILD_NUMBER_VALUE="$(echo ${TRAVIS_BUILD_NUMBER})"
GOBUILD_LDFLAGS="\
        -X ${VERSION_VAR}=${VERSION_VALUE} \
        -X ${REV_VAR}=${REV_VALUE} \
        -X ${REV_URL_VAR}=${REV_URL_VALUE} \
        -X ${GENERATED_VAR}=${GENERATED_VALUE} \
        -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} \
        -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} \
        -X ${DOCKER_IMAGE_VAR}=${DOCKER_IMAGE_VER} \
"

main() {

  case "$1" in
    --test)
    __test
    ;;
    --init)
    __init
    ;;
    --clean)
    __clean
    ;;
    --docs-deploy)
    __docs_deploy
    ;;
    --build)
    __build
    ;;
    --host-build)
    __host_build
    ;;
    --docker_deploy)
    __docker_deploy
    ;;
    *)
    echo "Error: Invalid argument '$1'" >&2
    exit 1
    ;;
  esac

}

__test() {

   cd ${ROOT}
   go test -v ./tests/scripts
   go test -v ./tests/workflow
}

__init() {

    mkdir -p ${TMP_DIR}
    cd ${ROOT}
    go mod vendor
    cp ${ROOT}/conf/config.dev.json ${ROOT}/conf/config.json
}

__clean() {

    rm -rf ${ROOT}/vendor
    rm -rf ${TMP_DIR}
    rm -rf ${HOME}/${ARCHIVE}
}

__docs_deploy() {

    cd ${ROOT}/doc/themes/default

    npm install
    gulp

    cd ${ROOT}/doc
    hugo

    cd ${ROOT}/doc/public

    git init
    echo -e "Starting to documentation commit.\n"
    git config --global user.email "support@e154.ru"
    git config --global user.name "delta54"

    git remote add upstream "https://${GITHUB_OAUTH_TOKEN}@github.com/e154/smart-home.git"
    git fetch upstream
    git reset upstream/gh-pages

    rev=$(git rev-parse --short HEAD)

    git add -A .

    set +o errexit

    git commit -m "rebuild pages at ${rev}"
    git push -q upstream HEAD:gh-pages

    echo -e "Done documentation deploy.\n"

    set -o errexit
}

__build_pingmq() {

    cd ${TMP_DIR}

    echo ""
    echo "build command:"
    echo "xgo --out=pingmq --targets=linux/*,windows/*,darwin/* ${ROOT}/cmd/pingmq"
    echo ""

    xgo --out=pingmq --targets=linux/*,windows/*,darwin/* ${ROOT}/cmd/pingmq

    chmod +x ${ROOT}/bin/pingmq
    cp ${ROOT}/bin/pingmq ${TMP_DIR}

}

__build() {

    __build_pingmq

    # build
    cd ${TMP_DIR}

    BRANCH="$(git name-rev --name-only HEAD)"

    if [[ $BRANCH == *"tags/"* ]]; then
      BRANCH="master"
    fi

    echo "BRANCH ${BRANCH}"

    echo ""
    echo "build command:"
    echo "xgo --out=${EXEC} --branch=${BRANCH} --targets=linux/*,windows/*,darwin/* --ldflags='${GOBUILD_LDFLAGS}' ${ROOT}"
    echo ""

    xgo --out=${EXEC} --branch=${BRANCH} --targets=linux/*,windows/*,darwin/* --ldflags="${GOBUILD_LDFLAGS}" ${ROOT}

    mkdir -p ${TMP_DIR}/api/server/v1/docs/swagger
    cp ${ROOT}/api/server/v1/docs/swagger/swagger.yaml ${TMP_DIR}/api/server/v1/docs/swagger/

    cp -r ${ROOT}/conf ${TMP_DIR}
    cp -r ${ROOT}/data ${TMP_DIR}
    cp -r ${ROOT}/snapshots ${TMP_DIR}
    cp ${ROOT}/LICENSE ${TMP_DIR}
    cp ${ROOT}/README* ${TMP_DIR}
    cp ${ROOT}/contributors.txt ${TMP_DIR}
    cp ${ROOT}/bin/docker/Dockerfile ${TMP_DIR}

    cp ${ROOT}/bin/server ${TMP_DIR}
    cp ${ROOT}/bin/server-installer.sh ${TMP_DIR}
    chmod +x ${TMP_DIR}/data/scripts/ping.sh

    cd ${TMP_DIR}

    ls -l

    echo "tar: ${ARCHIVE} copy to ${HOME}"

    # create arch
    tar -zcf ${HOME}/${ARCHIVE} .
}

__host_build() {

    OUTPUT="server-linux-amd64"

    echo ""
    echo "build command:"
    echo "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '${GOBUILD_LDFLAGS}' -o ${OUTPUT}"
    echo ""

    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "${GOBUILD_LDFLAGS}" -o ${OUTPUT}
}

__docker_deploy() {

    cd ${TMP_DIR}

    ls -ll

    echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

    # build image
    docker build -f ${ROOT}/bin/docker/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE} .
    # set tag to builded image
    docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_VER}
    docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_LATEST}
    # push tagged image
    docker push ${DOCKER_IMAGE_VER}
    docker push ${DOCKER_IMAGE_LATEST}
}

__test() {
    cd ${ROOT}

    go test -v ./tests/api
    go test -v ./tests/scripts
    go test -v ./tests/workflow
}

main "$@"
