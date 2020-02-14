#!/bin/bash

# Allow for this script to be called from anywhere
PARENT_DIR=$(dirname "$0")

# Move to the root of the project
# shellcheck disable=SC2164
cd "${PARENT_DIR}/.."

# Get the latest release from the upstream
LATEST_RELEASE=$(curl -s https://api.github.com/repos/golangci/golangci-lint/releases/latest  | grep "browser_download.*.checksums.txt" | cut -d '"' -f 4 | cut -d '/' -f 8)

# Edit this list if we want other files to be checked
FILES_TO_CHECK=( .travis.yml README.md )

for i in "${FILES_TO_CHECK[@]}"
do
    # shellcheck disable=SC2002
    CURRENT_VERSION=$(cat "${PARENT_DIR}"/../"${i}" | grep golangci | grep curl | awk '{ print $NF }')
    if [[ -n ${CURRENT_VERSION} ]]; then
        sed -i '' -e "s/${CURRENT_VERSION}/${LATEST_RELEASE}/g" "${PARENT_DIR}/../${i}"
    fi
done
