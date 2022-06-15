#!/bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )";

function test_install_ubuntu_20.04 {
    cat <<EOF | run_in_docker ubuntu:20.04
apt update && apt install -y curl
bash /src/misc/script/install.sh
test -f ./frmd 
EOF
}

function test_install_alpine_3.16 {
    cat <<EOF | run_in_docker alpine:3.16
apk add curl
sh /src/misc/script/install.sh
test -f ./frmd 
EOF
}

function test_install_fedora_36 {
    cat <<EOF | run_in_docker fedora:36
yum install -y util-linux
bash /src/misc/script/install.sh
test -f ./frmd
EOF
}

function run_in_docker {
    local image=$1
    cat | docker run \
        -v "${SCRIPT_DIR}/../..:/src" \
        --workdir /tmp \
        -i --rm \
        ${image}
}