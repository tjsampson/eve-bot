CI_COMMIT_BRANCH ?= local
CI_COMMIT_SHORT_SHA ?= 000001
CI_PROJECT_ID ?= 0
CI_PIPELINE_IID ?= 0
GOPATH ?= ${HOME}/go
MODCACHE ?= ${GOPATH}/pkg/mod
BUILD_NUMBER := ${CI_PIPELINE_IID}
PATCH_VERSION := $(shell cat VERSION)
VERSION := ${PATCH_VERSION}.${BUILD_NUMBER}
DOCKER_UID = $(shell id -u)
DOCKER_GID = $(shell id -g)
CUR_DIR := $(shell pwd)
BUILD_IMAGE := unanet-docker.jfrog.io/golang
IMAGE_NAME := unanet-docker-int.jfrog.io/ops/eve-bot-v1
LABEL_PREFIX := com.unanet
IMAGE_LABELS := \
	--label "${LABEL_PREFIX}.git_commit_sha=${CI_COMMIT_SHORT_SHA}" \
	--label "${LABEL_PREFIX}.gitlab_project_id=${CI_PROJECT_ID}" \
	--label "${LABEL_PREFIX}.build_number=${BUILD_NUMBER}" \
	--label "${LABEL_PREFIX}.version=${VERSION}"

docker-scanner-exec = docker run --rm \
	-e SONAR_TOKEN=${SONARQUBE_TOKEN} \
	-e SONAR_HOST_URL=https://sonarqube.unanet.io \
	-v $(CUR_DIR):/usr/src \
	--user="${DOCKER_UID}:${DOCKER_GID}" \
	sonarsource/sonar-scanner-cli sonar-scanner -Dsonar.projectKey=eve-bot -Dsonar.exclusions=**/*_test.go,**/*mock*.go

docker-exec = docker run --rm \
	-e DOCKER_UID=${DOCKER_UID} \
	-e DOCKER_GID=${DOCKER_GID} \
	-v ${CUR_DIR}:/src \
	-v ${MODCACHE}:/go/pkg/mod \
	-v ${HOME}/.ssh/id_rsa:/home/unanet/.ssh/id_rsa \
	-w /src \
	${BUILD_IMAGE}

docker-helm-exec = docker run --rm --user ${DOCKER_UID}:${DOCKER_UID} \
	-v ${CUR_DIR}:/src \
	-w /src \
	alpine/helm

check-tag = !(git rev-parse -q --verify "refs/tags/v${PATCH_VERSION}" > /dev/null 2>&1) || \
	(echo "the version: ${PATCH_VERSION} has been released already" && exit 1)

.PHONY: build dist test check_version

check_version:
	@$(check-tag)

build: check_version
	docker pull ${BUILD_IMAGE}
	docker pull unanet-docker.jfrog.io/alpine-base
	mkdir -p bin
	$(docker-exec) go build -ldflags="-X 'gitlab.unanet.io/devops/eve-bot/internal/api/api.Version=${VERSION}'" \
		-o ./bin/eve-bot ./cmd/eve-bot/main.go
	docker build . -t ${IMAGE_NAME}:${PATCH_VERSION}
	$(docker-helm-exec) package --version ${PATCH_VERSION} --app-version ${VERSION} ./.helm

test:
	docker pull ${BUILD_IMAGE}
	$(docker-exec) go build ./...
	$(docker-exec) go test -tags !local ./...

dist: build
	docker push ${IMAGE_NAME}:${PATCH_VERSION}
	curl --fail -H "X-JFrog-Art-Api:${JFROG_API_KEY}" \
		-X PUT \
		https://unanet.jfrog.io/unanet/api/storage/docker-int-local/ops/eve-bot-v1/${PATCH_VERSION}\?properties=version=${VERSION}%7Cgitlab-build-properties.project-id=${CI_PROJECT_ID}%7Cgitlab-build-properties.git-sha=${CI_COMMIT_SHORT_SHA}%7Cgitlab-build-properties.git-branch=${CI_COMMIT_BRANCH}
	curl --fail -H "X-JFrog-Art-Api:${JFROG_API_KEY}" \
			-T eve-bot-${PATCH_VERSION}.tgz "https://unanet.jfrog.io/artifactory/helm-local/eve-bot/eve-bot-${PATCH_VERSION}.tgz"

scan:
	$(docker-scanner-exec)