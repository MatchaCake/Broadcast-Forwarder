#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

lint_fast() {
  # Run only fast linter
  if has golangci-lint; then
    golangci-lint run --fast $*
  else
    ./bin/golangci-lint run --fast $*
  fi

}

# shellcheck disable=SC2120
lint_all() {
  # Run all the linter on all packages
  if has golangci-lint; then
    golangci-lint run $*
  else
    ./bin/golangci-lint run $*
  fi
}

case "$1" in
all)
  lint_all
  exit
  ;;
*)
  shift
  lint_fast $@
  exit
  ;;
esac

EXIT_CODE=$?
exit ${EXIT_CODE}
