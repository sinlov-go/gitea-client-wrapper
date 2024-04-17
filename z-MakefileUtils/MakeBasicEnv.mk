# this file must use as base Makefile
# use as:
# include z-MakefileUtils/MakeBasicEnv.mk
## feature:
# - can fetch PLATFORM OS_BIT ENV_ROOT ENV_HOME_PATH ENV_NOW_TIME_FORMAT, from runner
# - can fetch ENV_DIST_VERSION ENV_DIST_MARK , from CI/CD or git
## task:
# make envHelp
# make envBasic
# make envDistBasic

# env PLATFORM OS_BIT ENV_ROOT ENV_HOME_PATH
ifeq ($(OS),Windows_NT)
  PLATFORM=Windows
#  OS_BIT?=${shell if /i "%PROCESSOR_ARCHITECTURE%"=="AMD64" (echo x86_64) ELSE echo x86}
  OS_BIT?=$(shell echo `uname -m`) # x86_64 arm64
  # do windows powershell
  ENV_ROOT?=$(shell pwd)
  ENV_HOME_PATH?=${shell echo %UserProfile%}
  # ENV_NOW_TIME_FORMAT = $(shell echo %Date:~0,4%%Date:~5,2%%Date:~8,2%)
  ENV_NOW_TIME_FORMAT=$(shell echo %Date:~0,4%-%Date:~5,2%-%Date:~8,2%-%time:~0,2%-%time:~3,2%-%time:~6,2%)
else
  OS_UNAME?=$(shell echo `uname`) # Linux Darwin
  OS_BIT?=$(shell echo `uname -m`) # x86_64 arm64
  ENV_ROOT?=$(shell pwd)
  ENV_HOME_PATH?=${HOME}
  # ENV_NOW_TIME_FORMAT = $(shell date -u '+%Y-%m-%d-%H-%M-%S')
  ENV_NOW_TIME_FORMAT = $(shell date '+%Y-%m-%d-%H-%M-%S')
  ifeq ($(shell uname),Darwin)
    PLATFORM=MacOS
    ifeq ($(shell echo ${OS_BIT}),arm64)
      PLATFORM=MacOS-Apple-Silicon
    else
      PLATFORM=MacOS-Intel
    endif
    # do MacOS

  else
    PLATFORM=Unix-Like
    # do unix
  endif
endif

# change mark from woodpecker-ci https://woodpecker-ci.org/docs/usage/environment
ifneq ($(strip $(CI_COMMIT_TAG)),)
$(info -> change ENV_DIST_MARK by CI_COMMIT_TAG)
    ENV_DIST_MARK=-tag.${CI_COMMIT_TAG}
else
    ifneq ($(strip $(CI_COMMIT_SHA)),)
$(info -> change ENV_DIST_MARK by CI_COMMIT_SHA)
    ENV_DIST_MARK=-${CI_COMMIT_SHA}
    endif
endif
# change mark from https://docs.drone.io/pipeline/environment/substitution/
ifneq ($(strip $(DRONE_TAG)),)
$(info -> change ENV_DIST_MARK by DRONE_TAG)
    ENV_DIST_MARK=-tag.${DRONE_TAG}
else
    ifneq ($(strip $(DRONE_COMMIT)),)
$(info -> change ENV_DIST_MARK by DRONE_COMMIT)
    ENV_DIST_MARK=-${DRONE_COMMIT}
    endif
endif

# change mark from github actions https://docs.github.com/actions/learn-github-actions/environment-variables
ifneq ($(strip $(GITHUB_SHA)),)
$(info -> change ENV_DIST_MARK by GITHUB_SHA)
    ENV_DIST_MARK=-${GITHUB_SHA}
endif

# if above CI not set ENV_DIST_MARK, use git
ifeq ($(strip $(ENV_DIST_MARK)),)
$(info -> change ENV_DIST_MARK by git)
    ENV_DIST_MARK=-$(strip $(shell git --no-pager rev-parse --short HEAD))
endif

# finally change by ENV_CI_DIST_MARK
ifneq ($(strip $(ENV_CI_DIST_MARK)),)
$(info -> change ENV_DIST_MARK by ENV_CI_DIST_MARK)
    ENV_DIST_MARK=-${ENV_CI_DIST_MARK}
endif

# finally change by env ENV_CI_DIST_VERSION
ifneq ($(strip $(ENV_CI_DIST_VERSION)),)
$(info -> change ENV_DIST_VERSION by ENV_CI_DIST_VERSION)
    ENV_DIST_VERSION=${ENV_CI_DIST_VERSION}
endif

envBasic:
	@echo ------- start show env basic---------
	@echo ""
	@echo "PLATFORM                                  ${PLATFORM}"
	@echo "OS_BIT                                    ${OS_BIT}"
	@echo "ROOT_NAME                                 ${ROOT_NAME}"
	@echo "ENV_ROOT                                  ${ENV_ROOT}"
	@echo "ENV_NOW_TIME_FORMAT                       ${ENV_NOW_TIME_FORMAT}"
	@echo "ENV_HOME_PATH                             ${ENV_HOME_PATH}"
	@echo ""
	@echo ------- end  show env basic ---------


envDistBasic:
	@echo "ENV_DIST_VERSION :                        ${ENV_DIST_VERSION}"
	@echo "ENV_DIST_MARK :                           ${ENV_DIST_MARK}"
	@echo ""

envHelp:
ifeq ($(OS),Windows_NT)
	@echo ""
	@echo "warning: other install make cli tools has bug"
	@echo " run will at make tools version 4.+"
	@echo "windows use this kit must install tools blow:"
	@echo "-> scoop install main/make"
	@echo ""
	@echo "https://scoop.sh/#/apps?q=busybox&s=0&d=1&o=true"
	@echo "-> scoop install main/busybox"
	@echo "and if want use shasum must install tools blow:"
	@echo "https://scoop.sh/#/apps?q=shasum&s=0&d=1&o=true"
	@echo "-> scoop install main/shasum"
	@echo ""
endif
	@echo "=> make envBasic            - print basic env of this project"
	@echo "=> make envDistBasic        - print dist env of this project"