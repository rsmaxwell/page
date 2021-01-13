#!/bin/bash

. ~/web-server-info

if [ -z "${SERVER_SSH_USERNAME}" ]; then
    echo "SERVER_SSH_USERNAME not specified"
    exit 1
fi

if [ -z "${SERVER_SSH_HOST}" ]; then
    echo "SERVER_SSH_HOST not specified"
    exit 1
fi

if [ -z "${SERVER_SSH_PORT}" ]; then
    echo "SERVER_SSH_PORT not specified"
    exit 1
fi

if [ -z "${SERVER_RUN_USERNAME}" ]; then
    echo "SERVER_RUN_USERNAME not specified"
    exit 1
fi

if [ -z "${SERVER_SOURCE_DIRECTORY}" ]; then
    echo "SERVER_SOURCE_DIRECTORY not specified"
    exit 1
fi

if [ -z "${SERVER_TARGET_DIRECTORY}" ]; then
    echo "SERVER_TARGET_DIRECTORY not specified"
    exit 1
fi


target_dir="${SERVER_TARGET_DIRECTORY}/pages"
mkdir -p ${target_dir}


source="${SERVER_SOURCE_DIRECTORY}/diary-*"
target="${SERVER_SSH_USERNAME}@${SERVER_SSH_HOST}:${target_dir}/"

array=()
array+=("--rsh" "ssh -p ${SERVER_SSH_PORT}")
array+=("--rsync-path=sudo -u ${SERVER_RUN_USERNAME} rsync")
array+=("--archive")
array+=("--verbose")
array+=("--compress")
array+=("--human-readable")
array+=("--delete")
array+=("--exclude=**/metadata")
array+=("--exclude=*.db")
array+=("${source}")
array+=("${target}")

delimiter="|"
options=""
for i in "${array[@]}"; do
    if [[ "${options}" == "" ]]; then
        options=${i}
    else
        options="${options}${delimiter}${i}"
    fi
done

OFS=${IFS}
IFS=${delimiter}
set -x
rsync ${options}
set +x
result=$?
IFS=${OFS}

if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result = ${result}"
    exit 1
fi
