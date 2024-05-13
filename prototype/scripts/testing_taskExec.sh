#!/bin/bash

flag=0

echo "start" > out.txt
sudo lmu -T 10 -t 1 >> out.txt &
BACK_PID=$!

while kill -0 $BACK_PID ; do
    if ! pgrep bpftrace > /dev/null ; then
        if [[ $flag = 0 ]] ; then
            echo "no bpf" >> out.txt
            flag=1
        fi
    else
        if [[ $flag = 1 ]] ; then
            flag=0
        fi
    fi
done

wait $BACK_PID 
