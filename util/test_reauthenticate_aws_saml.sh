#!/usr/bin/env bash
# Tests AWS CLI for valid credentials and reauthenticates if credentials are
# invalid or expired.
# This script also manages the creation of a Python virtual environment.

set -Eeuo pipefail

msg() {
  echo >&2 -e "${1-}"
}

if ! [ -x "$(command -v python3)" ]; then
  msg 'Error: Python (python3) is not in the PATH, is it installed?'
  exit 1
fi

# The AWS profile that the credentials are saved to
AWS_PROFILE="saml"

if aws sts get-caller-identity --profile=$AWS_PROFILE > /dev/null 2>&1; then
    # Credentials are valid, no need to continue
    exit 0
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"

(
    cd $SCRIPT_DIR

    if [ ! -d $SCRIPT_DIR/venv ] ; then
        msg "Setting up the Python virtual environment"
        (
            python3 -m venv venv
            source ./venv/bin/activate
            python3 -m pip install -r ./requirements.txt
        )
        msg ""
    fi

    source ./venv/bin/activate

    msg "Enter your General Realm credentials:"
    python3 aws_saml_auth.py
)
