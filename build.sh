#!/bin/bash

export GOPATH=${WORKSPACE}
export PROJECT_DIR=${WORKSPACE}/src/github.com/rsmaxwell/page

export PAGE_ROOTDIR=${PROJECT_DIR}/testing/root
export DEBUG_DUMP_DIR=${PROJECT_DIR}/build/dumps

export repoUrl="https://server.rsmaxwell.co.uk/archiva"

gradle clean generate build
