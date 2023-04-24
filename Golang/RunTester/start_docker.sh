#!/bin/bash
cd $2
docker compose up -d
sleep 2
docker stop $1-chrome-video
sleep 10
docker start $1-chrome-video
sleep 2
pytest test.py 1> errors.txt 0> errors.txt

date_time=`date "+%d-%m-%Y-%H-%M-%S"`
docker stop $1-chrome-video
mv chrome.mp4 ./chrome-$1-video.mp4
#mv chrome.mp4 ./chrome-$1-$date_time-video.mp4
#echo chrome-$1-$date_time-video.mp4 >> $2/logs.txt


