#!/bin/bash
printf "\n" >> $1/.env
echo "FILE=$2" >> $1/.env
echo "PATHFOLDER=$1" >> $1/.env
echo "URL=$3" >> $1/.env
echo "USER1=$4" >> $1/.env

echo "PORT1=$5" >> $1/.env
echo "PORT2=$6" >> $1/.env
echo "PORT3=$7" >> $1/.env




