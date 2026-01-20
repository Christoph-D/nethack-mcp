#!/bin/bash

set -e

if [ ! -x ./nethack-llm/src/nethack ]; then
  echo "Building NetHack..."
  cd nethack-llm/sys/unix
  sh setup.sh hints/linux
  cd ../..
  make
  cd ..
fi

NETHACK_TMUX_SESSION=${NETHACK_TMUX_SESSION:-nethack}

NETHACK_DUMP_FILENAME=/tmp/${NETHACK_TMUX_SESSION}-map.json
export NETHACK_DUMP_FILENAME

tmux new-session -d -s "${NETHACK_TMUX_SESSION}" -n "${NETHACK_TMUX_SESSION}" bash -c '
while true; do
    NETHACKOPTIONS="role:wiz,race:hum,align:neu,gender:mal,!autopickup" \
      ./nethack-llm/src/nethack -u robot || exit
    echo "NetHack exited. Restarting in 2 seconds..."
    sleep 2 || exit
done
' \; \
  set-window-option -t "${NETHACK_TMUX_SESSION}" window-size manual \; \
  resize-pane -t 0 -x 80 -y 24 \; \
  attach -t "${NETHACK_TMUX_SESSION}"
