#!/bin/sh

kill_app() {
    pids=`pgrep "$1"`
    for pid in $pids
        do
            kill $pid
            echo "Killed process \"$1\" (pid: $pid)"
        done
}


kill_app $1