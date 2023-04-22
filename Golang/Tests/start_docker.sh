#!/bin/bash
#$3 testFile
# $2 - user
cd $1
#$1 - это просто пусть к папке юзера
docker-compose up -d
sleep 2
docker stop $2-chrome-video
sleep 10
docker start $2-chrome-video
sleep 2
pytest $3 1> errors.txt 0> errors.txt

date_time=`date "+%d-%m-%Y-%H-%M-%S"`
filename = "$2-$date_time-video.mp4"
docker stop $2-chrome-video && mv chrome.mp4 ./$filename

echo $filename


