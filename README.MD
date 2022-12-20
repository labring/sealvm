### how to using it in your project

- install vm tools
```shell
sealvm install 
```

- init vms 
```shell
sealvm run --nodes=1 
```

- init vms add dev 
```shell
 sealvm run --nodes=1 --dev --dev-mounts /Users/cuisongliu/Workspaces/go/src/github.com/labring/sealos:/root/go/src/github.com/labring/sealos 
```

- reset vms
```shell
sealvm reset
```


- install gcc
```shell
sealvm exec -r golang "sudo apt-get update"
sealvm exec -r golang "sudo apt-get install -y make gcc-aarch64-linux-gnu gcc-x86-64-linux-gnu"
```

- install golang
```shell
sealvm exec -r golang "bash golang-init"
sealvm exec -r golang "source /etc/profile.d/golang.sh && echo \$PATH"
sealvm exec -r golang "cd /root/go/src/github.com/labring/sealos && git config --global --add safe.directory /root/go/src/github.com/labring/sealos &&  source /etc/profile.d/golang.sh &&  make build"
```

- install k8s
```shell
sealvm exec -r golang "scp -o \"StrictHostKeyChecking no\" /root/go/src/github.com/labring/sealos/bin/linux_arm64/sealos root@192.168.64.55:/usr/bin/"
sealvm exec --hostnames  default-node-0 "sealos version" 
sealvm exec --hostnames  default-node-0 "sealos run labring/kubernetes:v1.25.0 --single" 
```
