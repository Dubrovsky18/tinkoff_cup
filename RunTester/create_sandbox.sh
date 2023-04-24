#!/bin/bash

if [ -d $1 ];
    then 
        echo 1
    else
        mkdir $1
fi

cp docker-compose.yml ./$1
if [ -f ./$1/.env ]
    then
    rm ./$1/.env
fi

cp .env ./$1/.env
cd $1

echo "PORT=$2" >> .env
echo "TEAM=$1" >> .env
echo "URL=$3" >> .env




