#!/bin/bash

sudo zypper install emacs
sudo rpm -q emacs
echo "install emacs emacs-27.2-150400.3.6.1.x86_64"

echo "sleep 10m in"
echo "Run the prototype..."
sleep 10m
echo "sleep 10m out"

file="packagesUsed.txt"

while read -r line; do
    package_file=$(sudo rpm -ql $line | head -n 1)
    head -c0 $package_file
    echo "open $package_file from package $line"
    sleep 1
done <$file
