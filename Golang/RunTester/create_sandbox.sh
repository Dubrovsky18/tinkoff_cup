#!/bin/bash

cp RunTester/docker-compose.yml $2

if [ -f ./$2/.env ]
    then
    rm ./$2/.env
fi

cp RunTester/.env $2/.env

echo "PORT=$4" >> $2/.env
echo "USER1=$1" >> $2/.env
echo "PATH=$2" >> $2/.env
echo "URL=$3" >> $2/.env




