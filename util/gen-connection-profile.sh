#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"
source "${SCRIPT_DIR}/util-common.bash"

if ! [ -x "$(command -v python3)" ]; then
  msg 'Error: Python (python3) is not in the PATH, is it installed?'
  exit 1
fi

test_setup_venv

source "${SCRIPT_DIR}/venv/bin/activate"
python3 "${SCRIPT_DIR}/gen-connection-profile.py" "$@"
