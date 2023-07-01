# SealVM 使用文档

SealVM是一个强大的虚拟机管理工具，主要用于方便地开发和测试集成的SealOS。以下是一些关键指令的说明和使用方法。

## VM管理命令

### 1. 运行(run)

该命令用于运行云原生的虚拟机节点。使用格式如下：

```shell
sealvm run --nodes=node:2,master:1
```

此命令将运行2个node角色和1个master角色的虚拟机。**角色可以自行定义**，输入要求 <role>:<count>

### 2. 重置(reset)

该命令用于重置虚拟机。使用格式如下：

```
sealvm reset
```

### 3. 检查(inspect)

该命令用于检查虚拟机节点的状态和配置。使用格式如下：

```
sealvm inspect <节点名称>
```

### 4. 列表(list)

该命令用于列出当前管理的所有虚拟机节点。使用格式如下：

```
sealvm list
```

## 远程操作命令

### 1. 操作(action)

该命令用于远程执行特定的操作。具体的操作参数需要根据具体情况填写。使用格式如下：

```
sealvm action -f action.yaml --debug
```

## 系统管理命令

### 安装(install) 新版本废弃(v0.2.0)

该命令用于安装虚拟机相关的工具。使用格式如下：

```
sealvm install <工具参数>
```

### 模板管理命令

#### 1. 设置默认模板(default)

此命令用于设置默认的虚拟机模板，可以同时指定多个角色。如：

```
sealvm template default -r node -r master
```

此命令将node和master设置为默认的角色。

#### 2. 获取模板(get)

此命令用于获取指定角色的模板。如：

```
sealvm template get node
```

此命令将获取node角色的模板。

#### 3. 列出模板(list)

此命令用于列出所有的模板。使用方法如下：

```
sealvm template list
```

#### 4. 设置模板(set)

此命令用于设置模板，具体的设置参数需要根据具体情况填写。如：

```
sealvm template set <模板参数>
```

#### 5. 重置模板(reset)

此命令用于重置模板。使用方法如下：

```
sealvm template reset
```

### 值管理命令

#### 1. 设置默认值(default)

此命令用于设置默认的值。使用方法如下：

```
sealvm values default
```

#### 2. 列出值(list)

此命令用于列出所有的值。使用方法如下：

```
sealvm values list
```

#### 3. 设置值(set)

此命令用于设置值，具体的设置参数需要根据具体情况填写。如：

```
sealvm values set <值参数>
```

### 配置管理命令

#### 1. 获取配置(get)

此命令用于获取指定的配置。如：

```
sealvm config get default_image
```

此命令将获取default_image的配置。

#### 2. 列出配置(list)

此命令用于列出所有的配置。使用方法如下：

```
sealvm config list
```

#### 3. 设置配置(set)

此命令用于设置配置，具体的设置参数需要根据具体情况填写。如：

```
sealvm config set default_image release:22.04
```

此命令将设置default_image的配置为release:22.04。

## 其他命令

### 1. 自动完成(completion)

该命令用于为指定的shell生成自动完成脚本。使用格式如下：

```
sealvm completion <shell参数>
```

对于每一个命令，如果需要更多的信息，可以使用 "sealvm <命令> --help" 这样的格式来获取更多帮助信息。希望这份文档能帮助你更好地使用SealVM。

