.PHONY: test check clean build dist all
#TOP_DIR := $(shell pwd)
# can change by env:ENV_CI_DIST_VERSION use and change by env:ENV_CI_DIST_MARK by CI
ENV_DIST_VERSION =v0.1.2
ENV_DIST_MARK=

ROOT_NAME ?=gitea-client-wrapper

## MakeDocker.mk settings start
ROOT_OWNER ?=sinlov-go
ROOT_PARENT_SWITCH_TAG =1.20.14
# for image local build
INFO_TEST_BUILD_DOCKER_PARENT_IMAGE =golang
# for image running
INFO_BUILD_DOCKER_FROM_IMAGE =alpine:3.17
INFO_BUILD_DOCKER_FILE =Dockerfile
INFO_TEST_BUILD_DOCKER_FILE =build.dockerfile
## MakeDocker.mk settings end
## MakeDockerCompose.mk settings start
INFO_DOCKER_COMPOSE_DEFAULT_FILE ?=docker-compose.yml
## MakeDockerCompose.mk settings end

## run info start
ENV_RUN_INFO_HELP_ARGS= -h
ENV_RUN_INFO_ARGS=
## run info end

## build dist env start
# change to other build entrance
ENV_ROOT_BUILD_ENTRANCE =cmd/gitea-client-wrapper/main.go
ENV_ROOT_BUILD_BIN_NAME =${ROOT_NAME}
ENV_ROOT_BUILD_PATH =build
ENV_ROOT_BUILD_BIN_PATH =${ENV_ROOT_BUILD_PATH}/${ENV_ROOT_BUILD_BIN_NAME}
ENV_ROOT_LOG_PATH =logs/
# linux windows darwin  list as: go tool dist list
ENV_DIST_GO_OS =linux
# amd64 386
ENV_DIST_GO_ARCH =amd64
# mark for dist and tag helper
ENV_ROOT_MANIFEST_PKG_JSON ?=package.json
ENV_ROOT_MAKE_FILE ?=Makefile
ENV_ROOT_CHANGELOG_PATH ?=CHANGELOG.md
## build dist env end

## go test MakeGoTest.mk start
# ignore used not matching mode
# set ignore of test case like grep -v -E "vendor|go_fatal_error" to ignore vendor and go_fatal_error package
ENV_ROOT_TEST_INVERT_MATCH?="vendor|go_fatal_error|robotn|shirou"
ifeq ($(OS),Windows_NT)
ENV_ROOT_TEST_LIST?=./...
else
ENV_ROOT_TEST_LIST?=$$(go list ./... | grep -v -E ${ENV_ROOT_TEST_INVERT_MATCH})
endif
# test max time
ENV_ROOT_TEST_MAX_TIME:=1m
## go test MakeGoTest.mk end

include z-MakefileUtils/MakeBasicEnv.mk
include z-MakefileUtils/MakeDistTools.mk
include z-MakefileUtils/MakeGoList.mk
include z-MakefileUtils/MakeGoMod.mk
include z-MakefileUtils/MakeGoTest.mk
include z-MakefileUtils/MakeGoTestIntegration.mk
include z-MakefileUtils/MakeGoDist.mk
# include MakeDockerRun.mk for docker run
include z-MakefileUtils/MakeDocker.mk
include z-MakefileUtils/MakeDockerCompose.mk

all: env

env: distEnv
	@echo "== project env info start =="
	@echo ""
	@echo "test info"
	@echo "ENV_ROOT_TEST_LIST                        ${ENV_ROOT_TEST_LIST}"
	@echo ""
	@echo "ROOT_NAME                                 ${ROOT_NAME}"
	@echo "ENV_DIST_VERSION                          ${ENV_DIST_VERSION}"
	@echo "ENV_ROOT_CHANGELOG_PATH                   ${ENV_ROOT_CHANGELOG_PATH}"
	@echo ""
	@echo "ENV_ROOT_BUILD_ENTRANCE                   ${ENV_ROOT_BUILD_ENTRANCE}"
	@echo "ENV_ROOT_BUILD_PATH                       ${ENV_ROOT_BUILD_PATH}"
ifeq ($(OS),Windows_NT)
	@echo "ENV_ROOT_BUILD_BIN_PATH                   $(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe"
else
	@echo "ENV_ROOT_BUILD_BIN_PATH                   ${ENV_ROOT_BUILD_BIN_PATH}"
endif
	@echo "ENV_DIST_GO_OS                            ${ENV_DIST_GO_OS}"
	@echo "ENV_DIST_GO_ARCH                          ${ENV_DIST_GO_ARCH}"
	@echo ""
	@echo "ENV_DIST_MARK                             ${ENV_DIST_MARK}"
	@echo "== project env info end =="

cleanBuild:
	@$(RM) -r ${ENV_ROOT_BUILD_PATH}
	@echo "~> finish clean path: ${ENV_ROOT_BUILD_PATH}"

cleanLog:
	@$(RM) -r ${ENV_ROOT_LOG_PATH}
	@echo "~> finish clean path: ${ENV_ROOT_LOG_PATH}"

cleanTest:
	@$(RM) coverage.txt
	@$(RM) coverage.out
	@$(RM) profile.txt

cleanTestData:
	$(info -> notes: remove folder [ testdata ] unable to match subdirectories)
	@$(RM) -r **/testdata
	@$(RM) -r **/**/testdata
	@$(RM) -r **/**/**/testdata
	@$(RM) -r **/**/**/**/testdata
	@$(RM) -r **/**/**/**/**/testdata
	@$(RM) -r **/**/**/**/**/**/testdata
	$(info -> finish clean folder [ testdata ])

clean: cleanTest cleanBuild cleanLog
	@echo "~> clean finish"

cleanAll: clean
	@echo "~> clean all finish"

init:
	@echo "~> start init this project"
	@echo "-> check version"
	go version
	@echo "-> check env golang"
	go env
	@echo "~> you can use [ make help ] see more task"
	-go mod verify

dep: modVerify modDownload modTidy
	@echo "-> just check depends below"

style: modTidy modVerify modFmt modLintRun

ci: modTidy modVerify modFmt modVet modLintRun test

ciTestBenchmark: modTidy modVerify testBenchmark

ciCoverageShow: modTidy modVerify modVet testCoverage testCoverageShow

ciAll: ci ciTestBenchmark ciCoverageShow

buildMain:
	@echo "-> start build local OS: ${PLATFORM} ${OS_BIT}"
ifeq ($(OS),Windows_NT)
	@go build -o $(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe ${ENV_ROOT_BUILD_ENTRANCE}
	@echo "-> finish build out path: $(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe"
else
	@go build -o ${ENV_ROOT_BUILD_BIN_PATH} ${ENV_ROOT_BUILD_ENTRANCE}
	@echo "-> finish build out path: ${ENV_ROOT_BUILD_BIN_PATH}"
endif

dev: export CI_DEBUG=true
dev: cleanBuild buildMain
ifeq ($(OS),Windows_NT)
	$(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe ${ENV_RUN_INFO_ARGS}
else
	${ENV_ROOT_BUILD_BIN_PATH} ${ENV_RUN_INFO_ARGS}
endif

runHelp: export CI_DEBUG=false
runHelp:
	go run -v ${ENV_ROOT_BUILD_ENTRANCE} ${ENV_RUN_INFO_HELP_ARGS}

helpProjectRoot:
	@echo "Help: Project root Makefile"
ifeq ($(OS),Windows_NT)
	@echo ""
	@echo "warning: other install make cli tools has bug"
	@echo " run will at make tools version 4.+"
	@echo "windows use this kit must install tools blow:"
	@echo "-> scoop install main/make"
	@echo ""
endif
	@echo "-- now build name: ${ROOT_NAME} version: ${ENV_DIST_VERSION}"
	@echo "-- distTestOS or distReleaseOS will out abi as: ${ENV_DIST_GO_OS} ${ENV_DIST_GO_ARCH} --"
	@echo ""
	@echo "~> make env                 - print env of this project"
	@echo "~> make init                - check base env of this project"
	@echo "~> make dep                 - check and install by go mod"
	@echo "~> make clean               - remove build binary file, log files, and testdata"
	@echo "~> make test                - run test case ignore --invert-match by config"
	@echo "~> make testCoverage        - run test coverage case ignore --invert-match by config"
	@echo "~> make testCoverageBrowser - see coverage at browser --invert-match by config"
	@echo "~> make testBenchmark       - run go test benchmark case all"
	@echo "~> make ci                  - run CI tools tasks"
	@echo "~> make ciTestBenchmark     - run CI tasks as test benchmark"
	@echo "~> make ciCoverageShow      - run CI tasks as test coverage and show"
	@echo "~> make ciAll               - run CI tasks all"
	@echo "~> make style               - run local code fmt and style check"
	@echo "~> make devHelp             - run as develop mode see help with ${ENV_RUN_INFO_HELP_ARGS}"
	@echo "~> make dev                 - run as develop mode"
	@echo "~> make runHelp             - run use ${ENV_RUN_INFO_HELP_ARGS}"

help: helpGoMod helpGoTest helpGoDist helpDocker helpProjectRoot
	@echo ""
	@echo "-- more info see Makefile include: MakeGoMod.mk MakeGoTest.mk MakeGoTestIntegration.mk MakeDocker.mk --"
