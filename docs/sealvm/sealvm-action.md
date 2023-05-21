# SealVM Action 使用文档

SealVM的 `Action` 是一种强大的操作，允许用户在虚拟机上执行一系列自定义任务。这些任务可以在一个或多个虚拟机上执行，支持挂载、卸载、执行命令、复制文件和写入文件内容等多种操作。`Action` 的配置需要一个YAML文件，以下是配置文件的各个部分的解释和使用方法。

## Action配置文件示例

```
apiVersion: virtual-machine.sealos.io/v1
kind: Action
spec:
  data:
    - mount:
        source: /Users/cuisongliu/Workspaces/go/src/github.com/labring/sealos
        target: /root/go/src/github.com/labring/sealos
    - umount: /root/go/src/github.com/labring/sealos
    - exec: |-
        ls -l /
        echo "ddff"
    - copy:
        source: /Users/cuisongliu/Workspaces/go/src/github.com/labring/sealos/README.md
        target: /root/README.md
    - exec: |-
        ls /root
    - copyContent:
        target: /root/sealos.sh
        content: |-
          write code
          dfff
    - exec: |-
        cat /root/sealos.sh
  ons:
    - role: master
    - role: worker
      indexes:
        - 0
        - 1
```

## 配置文件解释

- `apiVersion`：指定API的版本，固定为 `virtual-machine.sealos.io/v1`。
- `kind`：资源类型，固定为 `Action`。
- `spec`：具体的任务规格，由 `data` 和 `ons` 两部分组成。

    - `data`：一系列任务的列表。每个任务可以是以下五种类型之一：

        - `mount`：挂载一个目录或文件。需要提供源路径（`source`）和目标路径（`target`）。注意路径必须为绝对路径。
        - `umount`：卸载一个目录或文件。需要提供目标路径。
        - `exec`：在虚拟机上执行一系列命令。命令需要以字符串的形式给出，多条命令可以用换行符隔开。
        - `copy`：将一个文件从源路径复制到目标路径。需要提供源路径和目标路径。
        - `copyContent`：创建一个新文件，并写入指定的内容。需要提供目标路径和内容。

    - `ons`：指定任务要在哪些虚拟机上执行。每个虚拟机可以通过角色（`role`）和索引（`indexes`）来指定。如果不指定索引，则任务将在该角色的所有虚拟机上执行。

## 如何使用

1. 创建一个 `Action` 配置文件，按照上述格式编写你需要的任务。
2. 使用 `sealvm action -f <配置文件路径> --debug` 命令来执行 `Action`。`--debug` 参数是可选的，如果加

上，SealVM会打印更多的调试信息。

以上是SealVM Action的使用方法，希望能够帮助你更好地使用SealVM进行虚拟机管理。
