#!/bin/bash

while true; do
    nethack -u robot || exit
    echo "NetHack exited. Restarting in 2 seconds..."
    sleep 2
done
