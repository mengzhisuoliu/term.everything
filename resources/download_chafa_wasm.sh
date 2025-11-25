#!/bin/bash
VERSION="0.3.1"

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

OUTPUT_PATH="$SCRIPT_DIR/../framebuffertoansi/resources/chafa.wasm"
wget https://unpkg.com/chafa-wasm@$VERSION/dist/chafa.wasm -O "$OUTPUT_PATH"