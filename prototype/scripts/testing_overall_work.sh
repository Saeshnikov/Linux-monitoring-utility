#!/bin/bash

total_number_of_packages=$(rpm -qa | wc -l)
echo "The total number of packages on the system: $total_number_of_packages"
used_number_of_packages=$(( total_number_of_packages*80/100 ))
#echo $used_number_of_packages
array_packages=($(rpm -qa))


echo "sleep 10s in"
echo "Run the prototype..."
sleep 10
echo "sleep 10s out"

for ((i=0; i < used_number_of_packages; i++))
do
    number_of_random_package=$(shuf -i 0-$((total_number_of_packages-1)) -n1)
    #echo $number_of_random_package

    package_file=$(rpm -ql ${array_packages[number_of_random_package]} | head -n 1)
    head -c0 $package_file
    echo "open $package_file from package ${array_packages[number_of_random_package]}"
    sleep 1

done
