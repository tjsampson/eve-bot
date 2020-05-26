CI_COMMIT_BRANCH ?= local
CI_COMMIT_SHORT_SHA ?= 000001
CI_PROJECT_ID ?= 0
CI_PIPELINE_IID ?= 0
GOPATH ?= ${HOME}/go
MODCACHE ?= ${GOPATH}/pkg/mod

SONARQUBE_TOKEN := ${SONARQUBE_TOKEN}

VERSION_MAJOR := 0
VERSION_MINOR := 1
VERSION_PATCH := 0
BUILD_NUMBER := ${CI_PIPELINE_IID}
PATCH_VERSION := ${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}
VERSION := ${PATCH_VERSION}.${BUILD_NUMBER}

DOCKER_UID = $(shell id -u)
DOCKER_GID = $(shell id -g)

CUR_DIR := $(shell pwd)

BUILD_IMAGE := unanet-docker.jfrog.io/golang
IMAGE_NAME := unanet-docker.jfrog.io/eve-bot
IMAGE_DIGEST = $(shell docker inspect -f '{{index .RepoDigests 0}}' ${IMAGE_NAME}:${PATCH_VERSION})

LABEL_PREFIX := com.unanet
IMAGE_LABELS := \
	--label "${LABEL_PREFIX}.git_commit_sha=${CI_COMMIT_SHORT_SHA}" \
	--label "${LABEL_PREFIX}.gitlab_project_id=${CI_PROJECT_ID}" \
	--label "${LABEL_PREFIX}.build_number=${BUILD_NUMBER}" \
	--label "${LABEL_PREFIX}.version=${VERSION}"

docker-exec = docker run --rm \
	-e DOCKER_UID=${DOCKER_UID} \
	-e DOCKER_GID=${DOCKER_GID} \
	-v ${CUR_DIR}:/src \
	-v ${MODCACHE}:/go/pkg/mod \
	-v ${HOME}/.ssh/id_rsa:/home/unanet/.ssh/id_rsa \
	-w /src \
	${BUILD_IMAGE}

.PHONY: build dist test

build:
	docker pull ${BUILD_IMAGE}
	docker pull unanet-docker.jfrog.io/alpine-base
	mkdir -p bin
	$(docker-exec) go build -ldflags="-X 'gitlab.unanet.io/devops/eve/pkg/mux.Version=${VERSION}'" \
		-o ./bin/eve-bot ./cmd/eve-bot/main.go
	docker build . -t ${IMAGE_NAME}:${PATCH_VERSION}

test:
	docker pull ${BUILD_IMAGE}
	$(docker-exec) go build ./...
	$(docker-exec) go test -tags !local ./...

dist: build
	docker push ${IMAGE_NAME}:${PATCH_VERSION}
	curl --fail -H "X-JFrog-Art-Api:${JFROG_API_KEY}" \
		-X PUT \
		https://unanet.jfrog.io/unanet/api/storage/docker-local/eve-bot/${PATCH_VERSION}\?properties=version=${VERSION}%7Cgitlab-build-properties.project-id=${CI_PROJECT_ID}%7Cgitlab-build-properties.git-sha=${CI_COMMIT_SHORT_SHA}%7Cgitlab-build-properties.git-branch=${CI_COMMIT_BRANCH}

deploy:
	kubectl apply -f .kube/ingress.yaml
	kubectl apply -f .kube/manifest.yaml
	kubectl set image deployment/eve-bot-v1 eve-bot-v1=${IMAGE_DIGEST} --record

proxy-bot: 
	ssh evebot -R 3000:localhost:3000 -Nf

scan:
	docker run -e SONAR_TOKEN=${SONARQUBE_TOKEN} -e SONAR_HOST_URL=https://sonarqube.unanet.io -it -v $(pwd):/usr/src sonarsource/sonar-scanner-cli