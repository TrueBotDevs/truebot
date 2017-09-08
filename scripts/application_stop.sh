#!/bin/bash

service truebot stop
if [[ -e /home/ubuntu/truebot-2.0 ]]; then
  rm /home/ubuntu/truebot-2.0
fi
