/*
Copyright 2022 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tmpl

type Tpl string

// NodesTpl is the template of nodes.yaml
// Proxy is proxy addr exec, ex 192.168.64.1:7890
// NoProxy is no proxy addr exec, ex 192.168.0.0/16
// PrivateKey is private key file context
// PublicKey is public key file context
const NodesTpl Tpl = `write_files:
- content: |
    export https_proxy=http://{{ .Proxy }} http_proxy=http://{{ .Proxy }} all_proxy=socks5://{{ .Proxy }}
    export no_proxy=localhost,127.0.0.1,localaddress,.localdomain.com,apiserver.cluster.local,{{ .NoProxy }}
    echo -e "已开启代理"
  path: /usr/bin/proxy_on
  permissions: '0755'
- content: |
    unset http_proxy
    unset https_proxy
    unset ftp_proxy
    unset rsync_proxy
    echo -e "已关闭代理"
  path: /usr/bin/proxy_off
  permissions: '0755'
runcmd:
  - echo "deb [trusted=yes] https://apt.fury.io/labring/ /" | tee /etc/apt/sources.list.d/labring.list
  - echo "{{ .PublicKeyBase64 }}" | base64 -d >> /root/.ssh/authorized_keys
  - echo "{{ .PrivateKeyBase64 }}"| base64 -d > /root/.ssh/id_rsa
  - chmod 600 /root/.ssh/id_rsa
`

const GolangTpl Tpl = `write_files:
- content: |
    export https_proxy=http://{{ .Proxy }} http_proxy=http://{{ .Proxy }} all_proxy=socks5://{{ .Proxy }}
    export no_proxy=localhost,127.0.0.1,localaddress,.localdomain.com,apiserver.cluster.local,{{ .NoProxy }}
    echo -e "已开启代理"
  path: /usr/bin/proxy_on
  permissions: '0755'
- content: |
    unset http_proxy
    unset https_proxy
    unset ftp_proxy
    unset rsync_proxy
    echo -e "已关闭代理"
  path: /usr/bin/proxy_off
  permissions: '0755'
- content: |
    #!/bin/bash
    version=1.19.1
    arch={{ .ARCH }}
    rm -rf /root/go${version}.linux-${arch}.tar.gz
    wget https://studygolang.com/dl/golang/go${version}.linux-${arch}.tar.gz -O /root/go${version}.linux-${arch}.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -zxvf /root/go${version}.linux-${arch}.tar.gz
    echo "export PATH=\$PATH:/usr/local/go/bin" > /etc/profile.d/golang.sh
    chmod 0755 /etc/profile.d/golang.sh
    rm -rf /root/go${version}.linux-${arch}.tar.gz
    mkdir -p /root/go/src/github.com/labring /root/go/bin /root/go/pkg 
    go env -w GOPROXY="https://goproxy.io,direct"
    bash /etc/profile.d/golang.sh
  path: /usr/bin/golang-init
  permissions: '0755' 
- runcmd:
    - echo "{{ .PublicKeyBase64 }}" | base64 -d  >> /root/.ssh/authorized_keys
    - echo "{{ .PrivateKeyBase64 }}" | base64 -d > /root/.ssh/id_rsa
    - chmod 600 /root/.ssh/id_rsa
`
