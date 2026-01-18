#!/bin/bash

tmux new-session -d -s nethack -n nethack sh/run-nethack.sh robot \; \
  set-window-option -t nethack window-size manual \; \
  resize-pane -t 0 -x 80 -y 24 \; \
  attach -t nethack
