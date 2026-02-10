#!/bin/bash
mkdir -p debug/data{1,2}
for file in {test1,test2,test3}; do
head -c 1000000 /dev/urandom > debug/data1/$file
head -c 1000000 /dev/urandom > debug/data2/$file
done

