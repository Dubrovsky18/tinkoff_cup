#!/bin/bash
printf "\n" >> $1/.env
echo "FILE=$2" >> $1/.env
echo "PATHFOLDER=$1" >> $1/.env
echo "URL=https://www.tinkoff.ru/" >> $1/.env
echo "USER1=$4" >> $1/.env






