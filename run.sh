#!/bin/bash

tmux new-session -d -s nethack sh/run-nethack.sh \; \
  set-window-option -g window-size manual \; \
  resize-pane -t 0 -x 80 -y 24 \; \
  attach
