#!/bin/bash

while :
do
  ALIVE=`pidof node server.js | wc -l`
  if [ "${ALIVE}" -eq "1" ]; then
    sleep 1
  else
    echo -e '\033[1;33m>>> start CFCAPI...\033[0m'
    node server.js > /dev/null 2>&1 &
    while :
    do
      READY=`pidof node server.js | wc -l`
      if [ "${READY}" -eq "1" ]; then
        break;
      fi
      echo -n '.'
      sleep 0.1
    done
    echo -e '\033[1;33m OK\033[0m'
  fi
done
