#!/bin/bash

# NineLives
# Copyright (C) 2025 k-katfish
# Licensed under the NineLives License
# 2025-04-04

printHelp() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help        Show this help message"
    echo "  -b, --build       Build the project"
    echo "  -c, --clean       Clean the build directory"
    echo "  -t, --tiny        Build the tiny version"
}

BUILD=false
CLEAN=false
TINYVER=false

while [ -n "$1" ]; do
    case "$1" in
        -h|--help)
            printHelp
            exit 0
            ;;
        -b|--build)
            BUILD=true
            ;;
        -t|--tiny)
            TINYVER=true
            ;;
        -c|--clean)
            CLEAN=true
            ;;
        *)
            echo "Unknown option: $1"
            printHelp
            exit 1
            ;;
    esac
    shift
done

if ! $BUILD && ! $CLEAN; then
    echo "No options provided, assuming --build."
    BUILD=true
fi

if $CLEAN; then
    echo "Cleaning dist directory..."
    rm -rf dist/
    echo "Done."
fi

if $BUILD; then
    echo "Building project..."
    mkdir -p dist

    BUILD_CMD="go build -v"
    NLV="ninelives/internal/version"
    LDFLAGS="-X ${NLV}.Commit=$(git rev-parse HEAD) -X ${NLV}.BuildID=$(uuidgen)"

    BUILD_CMD="$BUILD_CMD -ldflags=\"${LDFLAGS}\""

    if $TINYVER; then
        BUILD_CMD="$BUILD_CMD -tags=tiny"
    fi

    BUILD_NLE="$BUILD_CMD -o dist/nle cmd/embed/main.go"
    BUILD_NLC="$BUILD_CMD -o dist/nlc cmd/client/main.go"

    echo "Building Embedder..."
    eval $BUILD_NLE && echo "Embedder built successfully." || echo "Error: Failed to build Embedder."
    
    echo "Building Client..."
    eval $BUILD_NLC && echo "Client built successfully." || echo "Error: Failed to build Client."

    echo "Build complete."
fi