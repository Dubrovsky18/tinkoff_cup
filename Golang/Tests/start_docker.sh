#!/bin/bash
#$3 testFile
# $2 - user
cd $1
#$1 - это просто пусть к папке юзера
docker compose up -d
sleep 15
docker stop $1-chrome-video
sleep 5
docker start $1-chrome-video
sleep 4
pytest $2 2> errors.txt 1> norm.txt

#date_time=date "+%d-%m-%Y-%H-%M-%S"
#filename="$1-video.mp4"
#docker stop $1-chrome-video
#mv chrome.mp4 ./$filename

pwd
echo $filename