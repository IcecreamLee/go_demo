#!/bin/bash
# 使用nohup命令后台运行crntab，并将输出写入crontab.log，运行后将pid写入crontab.pid
nohup ./crontab > crontab.log 2>&1 & echo $! > crontab.pid
