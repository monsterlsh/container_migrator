runc checkpoint --pre-dump --tcp-established --image-path /opt/container-migrator/client_repo/myredis/checkpoint0

```shell
imgDir="/opt/container-migrator/client_repo/myredis/checkpoint0"
pid=16144 
criu pre-dump  -t $pid  --track-mem -D $imgDir  
```