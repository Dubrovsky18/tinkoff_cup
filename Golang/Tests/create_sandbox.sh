#!/bin/bash

sudo cp docker-compose.yml $1/docker-compose.yml
if [ -f $1/.env ]
    then
    sudo rm $1/.env -y
fi

sudo cp .env $1/.env
cd $1

echo "FILE=$2" >> .env
echo "PATHFOLDER=$1" >> .env
echo "URL=$3" >> .env
echo "USER=$4" >> .env

echo "PORT1=$5" >> .env
echo "PORT2=$6" >> .env
echo "PORT3=$7" >> .env




