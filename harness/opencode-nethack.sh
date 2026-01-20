#!/bin/bash

set -e

cd "$(dirname "$0")"

PATH=$PATH:$(readlink -f ../go/bin)

if ! which nethack-ctl; then
  echo "Building nethack-ctl..."
  cd ..
  make
  cd "$(dirname "$0")"
  if ! which nethack-ctl; then
    echo "Failed to build nethack-ctl"
    exit 1
  fi
fi

NETHACK_TMUX_SESSION=${NETHACK_TMUX_SESSION:-nethack}
NETHACK_DUMP_FILENAME=/tmp/${NETHACK_TMUX_SESSION}-map.json

export NETHACK_TMUX_SESSION
export NETHACK_DUMP_FILENAME

opencode
