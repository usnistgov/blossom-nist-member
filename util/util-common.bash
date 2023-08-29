set -Eeuo pipefail

# The AWS profile that the credentials are saved to
export AWS_PROFILE="saml"

msg() {
  echo >&2 -e "${1-}"
}

# This will fail if SCRIPT_DIR
test_setup_venv() {
  local SOURCE_DIR
  SOURCE_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
  (
    cd "$SOURCE_DIR"

    if [ ! -d "$SOURCE_DIR"/venv ] ; then
        msg "Setting up the Python virtual environment"

        if ! [ -x "$(command -v python3)" ]; then
          msg 'Error: Python (python3) is not in the PATH, is it installed?'
          exit 1
        fi

        python3 -m venv venv
        source ./venv/bin/activate
        python3 -m pip install -r ./requirements.txt
        msg ""
    fi
  )
}