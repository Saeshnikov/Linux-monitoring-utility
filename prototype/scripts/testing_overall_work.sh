#!/bin/bash

cd /home/$USER/Desktop
mkdir checkPrototype
cd checkPrototype

mkdir NewFolderNew
cd NewFolderNew
touch file.txt
echo "create /Desktop/checkPrototype/NewFolderNew/file.txt"
cd /

echo "sleep 10m in"
echo "Run the prototype..."
sleep 10m
echo "sleep 10m out"

for i in {1..10}; do

    cd /home/$USER/Desktop/checkPrototype
    directory="NewFolder$i"
    mkdir $directory
    cd $directory

    touch open_file_$i
    cat /home/$USER/Desktop/checkPrototype/$directory/open_file_$i
    echo "open /Desktop/checkPrototype/$directory/open_file_$i"
    symlink="symlink$i"
    ln -s open_file_$i $symlink
    cat /home/$USER/Desktop/checkPrototype/$directory/symlink$i
    echo "open /Desktop/checkPrototype/$directory/symlink$i"
    touch not_open_file_$i
    #cat /home/$USER/Desktop/checkPrototype/$directory/not_open_file_$i
    echo "create /Desktop/checkPrototype/$directory/not_open_file_$i"

    echo "sleep 30m in"
    sleep 30m
    echo "sleep 30m out"
    #rm -R /home/$USER/Desktop/$directory
    cd /

done

cd /
