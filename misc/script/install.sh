#!/bin/sh
set -e

FORMIDABLE_RELEASES_URL="https://github.com/Bornholm/formidable/releases"
FORMIDABLE_DESTDIR="${FORMIDABLE_DESTDIR:-.}"
FORMIDABLE_FILE_BASENAME="frmd"
FORMIDABLE_PLATEFORM_SUFFIX="$(uname -s)_$(uname -m)"

function main {
    test -z "${FORMIDABLE_VERSION}" && FORMIDABLE_VERSION="$(curl -sfL -o /dev/null -w %{url_effective} "${FORMIDABLE_RELEASES_URL}/latest" |
        rev |
        cut -f1 -d'/'|
        rev)"

    # Check version variable initialization
    test -z "${FORMIDABLE_VERSION}" && {
        echo "Unable to get Formidable version !" >&2
        exit 1
    }

    test -z "${FORMIDABLE_TMPDIR}" && FORMIDABLE_TMPDIR="$(mktemp -d)"
    export TAR_FILE="${FORMIDABLE_TMPDIR}/${FORMIDABLE_FILE_BASENAME}_${FORMIDABLE_VERSION}_${FORMIDABLE_PLATEFORM_SUFFIX}.tar.gz"

    (
        cd "${FORMIDABLE_TMPDIR}"

        # Download Formidable
        echo "Downloading Formidable ${FORMIDABLE_VERSION}..."
        curl -sfLo "${TAR_FILE}" \
            "${FORMIDABLE_RELEASES_URL}/download/${FORMIDABLE_VERSION}/${FORMIDABLE_FILE_BASENAME}_${FORMIDABLE_VERSION}_${FORMIDABLE_PLATEFORM_SUFFIX}.tar.gz" ||
            ( echo  "Error while downloading Formidable !" >&2 && exit 1 )
        
        # Download checksums
        curl -sfLo "checksums.txt" "${FORMIDABLE_RELEASES_URL}/download/${FORMIDABLE_VERSION}/checksums.txt"
        
        echo "Verifying checksum..."
        check_sum ||
            ( echo  "Error while verifying checksums !" >&2 && exit 1 )
    )

    # Extracting archive files
    tar -xf "${TAR_FILE}" -C "${FORMIDABLE_TMPDIR}"

    # Moving downloaded binary to destination directory
    mv -f "${FORMIDABLE_TMPDIR}/${FORMIDABLE_FILE_BASENAME}" "${FORMIDABLE_DESTDIR}/"

    echo "You can now use '${FORMIDABLE_DESTDIR}/${FORMIDABLE_FILE_BASENAME}', enjoy !"
}

function check_sum {
    set -o pipefail
    cat checksums.txt | grep frmd_*_${FORMIDABLE_PLATEFORM_SUFFIX}.tar.gz | sha256sum -c
    set +o pipefail
}

main $@
