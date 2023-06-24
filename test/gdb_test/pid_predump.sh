pid=$(runc list | grep myredis | awk '{print $2}')
echo "running container pid is $pid"
criu pre-dump  -t $pid -v4  --track-mem -D /opt/container-migrator/client_repo/gdbmyredis/checkpoint -o ./log.txt