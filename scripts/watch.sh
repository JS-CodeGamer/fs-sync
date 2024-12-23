#!/bin/bash

command -v watchexec >/dev/null
if [ ! $? ]; then
  command -v cargo >/dev/null
  if [ ! $? ]; then
    printf '\e[31m error: watchexec required \e[0m \n'
    exit 1
  fi
  printf 'installing watchexec using cargo\n'
  cargo install watchexec
fi

if [ $# -ne 1 ]; then
  printf 'Usage: %s [command]\n' $0
  printf '  command:\n'
  printf '    run: runs run.sh\n'
  printf '    test: runs test.sh\n'
  exit 1
fi

if [ $1 = 'run' ]; then
  watchexec -c -re go,yaml --stop-signal SIGINT -- scripts/run.sh
elif [ $1 = 'test' ]; then
  watchexec -c -re go,yaml,sh -- scripts/test.sh
fi
