# Proxypool 健康检查（魔改）

> 此fork版本魔改自[原版本](https://github.com/Sansui233/proxypoolCheck)
> 在原先版本的基础上增加了`proxypool连接重试机制`和`自定义可用性检测网址`，以及一些bug修复（比如server_url自动添加后缀"/clash/proxies"的问题）

## 导航
- [简介](#简介)
- [安装和运行](#安装和运行)
- [配置](#配置)
- [添加自启](#添加自启)
- [声明](#声明)

## 简介

此项目为[proxypool](https://github.com/sansui233/proxypool)的代理池节点可用性检测部分，并提供检测后可用的代理。

此项目推荐在本地（即您家里）部署，或是的中国大陆服务器上运行，以提升代理池节点的实际可用比例。

在使用此项目之前，您应该有一个（或知道一个）可用的[proxypool](https://github.com/sansui233/proxypool)服务器。


## 安装和运行

二选一

### 1. 用构建好的（本fork版本没有构建好的，这里是原版本的链接）

从[releases](https://github.com/Sansui233/proxypoolCheck/releases)中下载

将下载的文件重命名为proxypoolcheck（可选）

不要忘了给文件添加755权限，否则无法运行

```
chmod +775 proxypoolcheck
```

您可以将config.yaml放在与proxypoolcheck文件同一文件夹内，或使用 -c 指定配置路径

```shell
./proxypoolCheck
# or
./proxypoolCheck -c /指定目录/config.yaml
```

### 2. 源代码编译/运行（推荐）

确保安装golang，然后下载源码
```sh
$ go get -u -v github.com/tagbug/proxypoolCheck
```

#### 先编译

```shell
$ go build main.go
```

> 生成的可执行文件为main.exe（windows）或main（linux）

**注意**：这条命令会根据你`当前的操作系统环境`编译成对应的可执行文件，如果你希望编译指定环境的可执行文件，请使用下面的命令

编译为linux可执行文件
```shell
$ cmd # 切换到cmd环境（如果你当前不是cmd环境的话）
$ SET CGO_ENABLED=1
$ SET GOOS=windows
$ SET GOARCH=amd64
$ go build main.go
```

编译为windows可执行文件
```shell
$ cmd # 切换到cmd环境（如果你当前不是cmd环境的话）
$ SET CGO_ENABLED=1
$ SET GOOS=windows
$ SET GOARCH=amd64
$ go build main.go
```

#### 直接运行

```shell script
$ go run main.go -c ./config/config.yaml
```

## 配置

基本的配置

```yaml
# proxypool远程服务器的地址，空白为http://127.0.0.1:8080
server_url:
  - https://example.proxypoolserver.com
  - https://example.proxypoolserver.com/clash/proxies?type=vmess


# 对于您的本地服务器
request: http   # http / https
domain:         # default: 127.0.0.1
port:           # default: 80

cron_interval: 15       # default: 15  minutes
show_remote_speed: true # default false

healthcheck_timout:     # default 5
healthcheck_connection: # default 100

speedtest:            # default false
speed_connection:     # default 5
speed_timeout:         # default 10
```

新增的配置

```yaml
server_max_retry: 5                           # 连接server_url的最大重试次数，default 1
healthcheck_test_site: https://t.me/tagbug    # 自定义可用性测试的网址，default https://youtube.com
```

需要修改的参数：

- `server_url`：远程服务器链接，可以使用筛选参数。支持多种来源
- `request`：要显示到网页的协议，默认 http，可选 https。
- `domain`：要显示到网页的域名，默认 127.0.0.1。
- `port`：要显示到网页上的端口，默认 80。如果本机有其他程序占用需要修改。

可选参数：

- `show_remote_speed`：显示远程速度，默认false，但建议改成true（因为作者写的就是true）
- `cron_interval`：工作间隔，默认15分钟
- `speedtest`：是否开启测速，默认关闭。开启测速会消耗大量流量。
- `speed_connection`：测速并发连接数，默认值为 5。
- `speed_timeout`：单个节点测速时间限制，默认值为 10，单位为秒。超过此时间限制的节点会测速失败
- `healthcheck_timeout`：单个节点健康检测时间限制，默认值为 5，单位为秒。超过此时间限制的节点为无效节点
- `healthcheck_connection`：节点健康检测并发连接数，默认值为 100。丢失大量可用节点时可大幅减少此项数值。
- `server_max_retry`：远程服务器连接失败的重试次数，默认值为 1
- `healthcheck_test_site`：自定义可用性测试的网址，默认值为 https://youtube.com

如果您的Web服务器端口与proxypoolCheck服务端口不同，应该将web服务器端口放在配置中，并且设置环境变量`PORT`以供proxypoolCheck服务。当您使用frp时，这将非常有帮助。

```shell
export PORT=ppcheckport
```
## 添加自启

此部分适用于Linux。

**配置 systemd 服务**

`vim /etc/systemd/system/proxypoolcheck.service` 填入下面内容：
```
[Unit]
Description=proxypoolcheck
After=network-online.target
 
[Service]
Type=simple
Restart=on-abort
ExecStart=/proxypoolcheck所在的目录/proxypoolcheck -c /指定配置文件目录/config.yaml
 
[Install]
WantedBy=default.target
```

**重载 systemd 服务**

```
systemctl daemon-reload
```

**启动proxypoolcheck服务**
```
systemctl start proxypoolcheck
```
执行`systemctl status proxypoolcheck`确认有以下信息

```
● proxypoolcheck.service - proxypoolcheck
     Loaded: loaded (/etc/systemd/system/proxypoolcheck.service; enabled; vendor preset: enabled)
     Active: active (running) since Sun 2021-03-21 14:53:55 UTC; 9s ago
```

**添加开机启动**
```
systemctl enable proxypoolcheck
```

**查询服务是否开机启动，enabled即开启自启**
```
systemctl is-enabled proxypoolcheck.service
```
**`reboot`重启后`systemctl status proxypoolcheck`看看是否正常，如果正常，您就可以给个star，然后关闭网页，尽情享受**


## 声明

本项目遵循 GNU General Public License v3.0 开源，在此基础上，所有使用本项目提供服务者都必须在网站首页保留指向本项目的链接

本项目仅限个人自己使用，**禁止使用本项目进行营利**和**做其他违法事情**，产生的一切后果本项目概不负责。