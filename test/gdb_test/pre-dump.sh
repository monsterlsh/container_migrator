#!/bin/bash
imgDir="/opt/container-migrator/client_repo/gdbmyredis/checkpoint"

cd ../../workloads/redis
rm data/*.rdb -rf
runc run myredis > /dev/null &
sleep 10
pid=$(runc list | grep myredis | awk '{print $2}')
# gdb criu 
echo "running container $pid"
criu pre-dump  -t $pid -v4  --track-mem -D /opt/container-migrator/client_repo/gdbmyredis/checkpoint -o ./log.txt