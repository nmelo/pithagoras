#!/usr/bin/env bash
set -e -o pipefail

. build-conf.sh

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

WITH_BUILDER=$2
WITH_BUILDER=${WITH_BUILDER:-1}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT" "$WITH_BUILDER" 
fi

echo
echo "==========================="
echo "Stamping the release with proper version"
./version-control.sh

echo "----------------------------"
echo "Packaging standalone release"
./package-standalone-release.sh "$WITH_BUILDER"

echo "------------------------------"
echo "Compressing standalone release"
./compress-standalone-release.sh

popd >/dev/null
