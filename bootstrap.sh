#!/bin/bash

pid=`ps -ef | grep "hertz_service" | grep -v grep | awk '{print $2}'`
echo "pid is $pid"
if [ ! -z $pid ]; then
  kill -9 $pid
  echo "kill $pid"
fi

time=$(date +%y%m%d%H%S)
echo $time
if mv /opt/pdm-plugin/output/logs/ /opt/pdm-plugin-persist/logs-$time; then
 echo 'mv success'
else
 echo 'mv failed'
 exit 1
fi
rm -rf output
tar -xvf pdm-plugin
cd output
./bootstrap.sh
netstat -tunlp | grep hertz_service
netstat -tunlp | grep hertz_service > /opt/pdm-plugin/last_pid
