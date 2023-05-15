## 构建本地镜像
### 编写 Dockerfile 将模块二作业编写的 httpserver 容器化
```bash
# 新建目录，为context做准备
$ mkdir /home/alfie/part3
$ cd /home/alfie/part3
# 获得此前httpserver代码
$ wget https://raw.githubusercontent.com/Alopex4/cloudnative/master/part2.golang/lesson5/main.go
# 查看dockerfile
$ cat Dockerfile 
FROM golang:alpine3.18 AS build
RUN mkdir /go/src/project
COPY main.go  /go/src/project
EXPOSE 8080
WORKDIR /go/src/project
RUN go mod init && go mod tidy && go build -o /bin/simpleHttpServer

FROM scratch
COPY --from=build /bin/simpleHttpServer /bin/simpleHttpServer
ENTRYPOINT ["/bin/simpleHttpServer"]
# 在当前context下创建名为 simple-http:latest 的镜像
$ docker build -t simple-http:v1 .
```
### 将镜像推送至 docker 官方镜像仓库
```bash
$ docker login -u "naronav" -p "mypassword" docker.io
$ docker tag simple-http:v1 naronav/simple-http:v1
$ docker push naronav/simple-http:v1
```
### 通过 docker 命令本地启动 httpserver
```bash
$ docker pull naronav/simple-http:v1
$ docker run -d -P naronav/simple-http:v1
e583123f0e1c75affc12fc7ebb70526cc42aedd09a262038bb0fb50714be076e
```
### 通过 nsenter 进入容器查看 IP 配置
```bash
# root权限执行，查看namespace的 net域 (其PID为9730)
# lsns -t net
        NS TYPE NPROCS   PID USER    NETNSID NSFS                           COMMAND
4026531992 net     178     1 root unassigned                                /usr/lib/systemd/systemd --switched-root --system --deserialize 17
4026532594 net       1  9730 root          2 /run/docker/netns/86821bbe65e2 /bin/simpleHttpServer

# 在PID为9730的网络ns下，执行 ip addr查看IP配置
# nsenter -t 9730 -n ip addr 
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
22: eth0@if23: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever

```