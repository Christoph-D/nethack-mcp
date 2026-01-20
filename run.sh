#!/bin/bash

tmux new-session -d -s nethack -n nethack bash -c '
while true; do
    NETHACKOPTIONS="role:wiz,race:hum,align:neu,gender:mal,!autopickup" \
      ./nethack-llm/src/nethack -u robot || exit
    echo "NetHack exited. Restarting in 2 seconds..."
    sleep 2 || exit
done
' \; \
  set-window-option -t nethack window-size manual \; \
  resize-pane -t 0 -x 80 -y 24 \; \
  attach -t nethack
