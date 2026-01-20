#!/bin/bash

while true; do
    NETHACKOPTIONS='role:wiz,race:hum,align:neu,gender:mal,!autopickup' \
      ./nethack-llm/src/nethack -u "$1" || exit
    echo "NetHack exited. Restarting in 2 seconds..."
    sleep 2
done
