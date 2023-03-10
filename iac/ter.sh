#!/usr/bin/env bash
# Wrapper for Terraform that checks for valid credentials

set -Eeuo pipefail

msg() {
  echo >&2 -e "${1-}"
}


if ! [ -x "$(command -v terraform)" ]; then
  msg 'Error: Python (terraform) is not in the PATH, is it installed?'
  exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"

(
    cd $SCRIPT_DIR

    ./../util/test_reauthenticate_aws_saml.sh

    AWS_PROFILE="saml" terraform "$@"
)
