#!/usr/bin/env bash
set -e -o pipefail

# Builds an entire electron-based GUI for release

# Implemented architectures:
#       darwin/amd64
#       windows/amd64
#       windows/386
#       linux/amd64
#
# By default builds all architectures.
# A single arch can be built by specifying it using gox's arch names

. build-conf.sh

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

WITH_BUILDER=$2
WITH_BUILDER=${WITH_BUILDER:-1}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT" "$WITH_BUILDER"
fi

if [ -e "$ELN_OUTPUT" ]; then
    rm -r "$ELN_OUTPUT"
fi

if [ "$WITH_BUILDER" = "1" ]; then
    if [ ! -z "$WIN64_ELN" ] && [ ! -z "$WIN32_ELN" ]; then
        npm run dist-win
    fi

    if [ ! -z "$LNX64_ELN" ]; then
        npm run dist-linux
    fi

    if [ ! -z "$OSX64_ELN" ]; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            npm run dist-mac
        elif [[ "$OSTYPE" == "linux"* ]]; then
            npm run pack-mac
        else
            echo "Can not run build script in $OSTYPE"
        fi
    fi

    pushd "$FINAL_OUTPUT" >/dev/null
    if [ -e "mac" ]; then
        pushd "mac" >/dev/null
        if [ -e "${PDT_NAME}-${APP_VERSION}.dmg" ]; then
            mv "${PDT_NAME}-${APP_VERSION}.dmg" "../${PKG_NAME}-${APP_VERSION}-gui-osx-x64.dmg"
        elif [ -e "${PDT_NAME}.app" ]; then
            tar czf "../${PKG_NAME}-${APP_VERSION}-gui-osx-x64.zip" --owner=0 --group=0 "${PDT_NAME}.app"
        fi
        popd >/dev/null
        rm -rf "mac"
    fi

    IMG="${PKG_NAME}-${APP_VERSION}-x86_64.AppImage"
    DEST_IMG="${PKG_NAME}-${APP_VERSION}-gui-linux-x64.AppImage"
    if [ -e $IMG ]; then
        mv "$IMG" "$DEST_IMG"
        chmod +x "$DEST_IMG"
    fi

    EXE="${PDT_NAME} Setup ${APP_VERSION}.exe"
    if [ -e "$EXE" ]; then
        mv "$EXE" "${PKG_NAME}-${APP_VERSION}-gui-win-setup.exe"
    fi

    # clean unpacked folders
    rm -rf *-unpacked

    popd >/dev/null

else
    GULP_PLATFORM=""
    if [ -n "$1" ]; then
        GOX_OSARCH="$1"
        case "$1" in
        "linux/amd64")
            GULP_PLATFORM="linux-x64"
            ;;
        "linux/arm")
            GULP_PLATFORM="linux-arm"
            ;;
        "windows/amd64")
            GULP_PLATFORM="win32-x64"
            ;;
        "windows/386")
            GULP_PLATFORM="win32-ia32"
            ;;
        "darwin/amd64")
            GULP_PLATFORM="darwin-x64"
            ;;
        esac
    fi

    if [ -n "$GULP_PLATFORM" ]; then
        gulp electron --platform "$GULP_PLATFORM"
    else
        gulp electron
    fi

    echo "--------------------------"
    echo "Packaging electron release"
    ./package-electron-release.sh

    echo "----------------------------"
    echo "Compressing electron release"
    ./compress-electron-release.sh
fi

popd >/dev/null
