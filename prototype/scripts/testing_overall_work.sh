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
sleep 10m
echo "sleep 10m out"

for i in {1..10}; do

    cd /home/$USER/Desktop/checkPrototype
    directory="NewFolder$i"
    mkdir $directory
    cd $directory

    touch $i
    cat /home/$USER/Desktop/checkPrototype/$directory/$i
    echo "open /Desktop/checkPrototype/$directory/$i"
    symlink="symlink$i"
    ln -s $i $symlink
    cat /home/$USER/Desktop/checkPrototype/$directory/symlink$i
    echo "open /Desktop/checkPrototype/$directory/symlink$i"
    touch $((i + i))
    #cat /home/$USER/Desktop/checkPrototype/$directory/$((i + i))
    echo "create /Desktop/checkPrototype/$directory/$((i + i))"

    echo "sleep 30m in"
    sleep 30m
    echo "sleep 30m out"
    #rm -R /home/$USER/Desktop/$directory
    cd /

done

cd /
