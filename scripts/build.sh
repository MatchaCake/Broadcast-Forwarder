#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

cd "$PROJECT_DIR" || exit 1

build() {
  echo_info "Build boardcast forwarder"
  GO111MODULE=on go build -v -i -o ./bin/boardcast-forwarder ./
  echo_info "Done. Here is target information"
  ls -lah -d "$PWD/bin/boardcast-forwarder"
}

show_help() {
  cat <<EOF
Build artifacts for Item Info project. Known recipes are:
   build        build main.go into ./bin/boardcast-forwarder,
                which is the final deployable target of this project.

If you are seeing this from "make build" command, then, keep in mine that those
recipes should be used in the form "make build.xxx", for example:

    make build.build
EOF
}

while :; do
  case $1 in
    -h | --help)
      show_help # Display a usage synopsis.
      exit
      ;;
    build)
      build
      exit
      ;;
    * | -?*)
      show_help
      break
      ;;
    *) # Default case: No more options.
      break ;;
  esac

  shift
done

cd "$WORKING_DIR" || exit 1
