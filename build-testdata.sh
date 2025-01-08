#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# TODO: the loop with the container in the middle is slow, let see if we can use an alternative
for filename in ${SCRIPT_DIR}/testdata-project/*.go; do

    name=$(basename ${filename} .go)
    docker run --rm \
        -v ${SCRIPT_DIR}/:/src \
        -w /src tinygo/tinygo:0.34.0 bash \
        -c "cd testdata-project && tinygo build --no-debug -target=wasm-unknown -o /tmp/tmp.wasm ${name}.go && cat /tmp/tmp.wasm" > \
        ${SCRIPT_DIR}/it/testdata/${name}.wasm

done
