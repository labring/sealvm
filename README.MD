### how to using it in your project

- install vm tools
```shell
sealvm install 
```

- init vms 
```shell
sealvm run --nodes=node:2,master:1
```


- run actions

```yaml
apiVersion: virtual-machine.sealos.io/v1
kind: Action
spec:
  data:
  - mount:
      source: /Users/cuisongliu/Workspaces/go/src/github.com/labring/sealos
      target: /root/go/src/github.com/labring/sealos
  - exec: |
     sudo apt-get -y update
     sudo apt-get install -y make gcc-aarch64-linux-gnu gcc-x86-64-linux-gnu
     sudo apt-get install -y ntpdate
     sudo ntpdate -s ntp1.aliyun.com
  - copy:
      source: /Users/cuisongliu/Workspaces/go/src/github.com/labring/sealos/README.md
      target: /root/README.md
  - copyContent: 
      content: |
        #!/bin/bash
        version=1.20.1
        arch=arm64
        rm -rf /root/go${version}.linux-${arch}.tar.gz
        wget https://studygolang.com/dl/golang/go${version}.linux-${arch}.tar.gz -O /root/go${version}.linux-${arch}.tar.gz
        rm -rf /usr/local/go && tar -C /usr/local -zxvf /root/go${version}.linux-${arch}.tar.gz
        echo "export PATH=\$PATH:/usr/local/go/bin" > /etc/profile.d/golang.sh
        chmod 0755 /etc/profile.d/golang.sh
        rm -rf /root/go${version}.linux-${arch}.tar.gz
        mkdir -p /root/go/src/github.com/labring /root/go/bin /root/go/pkg
        source /etc/profile.d/golang.sh
        go env -w GOPROXY="https://goproxy.io,direct"
      target: /root/golang-install.sh
  - exec: |
      bash /root/golang-install.sh
      source /etc/profile.d/golang.sh && echo $PATH
      git config --global --add safe.directory /root/go/src/github.com/labring/sealos
      cd /root/go/src/github.com/labring/sealos &&  source /etc/profile.d/golang.sh &&  make build
  - copyContent:
    content: |
      #!/bin/bash
      echo "deb [trusted=yes] https://apt.fury.io/labring/ /" | tee /etc/apt/sources.list.d/labring.list
      sudo apt-get update
      sudo apt-get install -y sealos
    target: /root/sealos-apt.sh
  ons:
  - role: master
```

```shell
sealvm action -f action.yaml
```

- reset vms

```shell
sealvm reset
```
