# 课后练习 4.1
## 用 Kubeadm 安装 Kubernetes 集群。

### 系统及软件信息
```c
/* hardware */
+------+---------+-----+------+
| HOST |   CPU   | RAM | DISK |
+------+---------+-----+------+
| m1   | 4 cores | 4GB | 30GB |
| n1   | 4 cores | 4GB | 30GB |
| n2   | 4 cores | 4GB | 30GB |
+------+---------+-----+------+

/* software */
+------+----------+----------------+---------------+--------------+
| HOST |   K8S    |      CRI       |      CNI      |     ETCD     |
+------+----------+----------------+---------------+--------------+
| m1   | v.1.22.2 | docker:20.10.7 | calico:3.23.5 | etcd:3.5.0-0 |
| n1   | v.1.22.2 | docker:20.10.7 | calico:3.23.5 | etcd:3.5.0-0 |
| n2   | v.1.22.2 | docker:20.10.7 | calico:3.23.5 | etcd:3.5.0-0 |
+------+----------+----------------+---------------+--------------+

/* system */
+-----------+---------------------------+------+---------------------+-------------------+--------------------+---------------+
|   Role    |            VM             | HOST |         OS          |      Kernel       |         IP         |    Gateway    |
+-----------+---------------------------+------+---------------------+-------------------+--------------------+---------------+
| ks-master | VMware Workstation 17 Pro | m1   | Ubuntu 20.04 Server | 5.4.0-159-generic | 192.168.100.120/24 | 192.162.100.2 |
| ks-node1  | VMware Workstation 17 Pro | n1   | Ubuntu 20.04 Server | 5.4.0-159-generic | 192.168.100.121/24 | 192.162.100.2 |
| ks-node2  | VMware Workstation 17 Pro | n2   | Ubuntu 20.04 Server | 5.4.0-159-generic | 192.162.100.122/24 | 192.168.100.2 |
+-----------+---------------------------+------+---------------------+-------------------+--------------------+---------------+
```

1. 安装前准备工作
```bash
# 更新源文件
cat <<EOF | sudo tee /etc/apt/sources.list
# use aliyun mirrors
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
EOF

# 删除cache清理 APT 软件包管理器缓存、自动清理 APT 软件包管理器缓存(删除过期软件包)、更新源
sudo apt clean && sudo apt autoclean && sudo apt update

# 追加主机名主机名解析
# 此前在安装的时候已修改了主机名，如果未设置主机名
# 可以通过 systemctl set-hostname 进行修改
sudo sed  -i '1s !^!192.168.100.120 m1\n192.168.100.121 n1\n192.168.100.122\n!'  /etc/hosts

# 关闭swap
sudo swapoff -a
sed -i '/swap/ s/^\(.*\)/#\1/' /etc/fstab

## 验证
free -h 
              total        used        free      shared  buff/cache   available
Mem:          3.8Gi       309Mi       2.7Gi       1.0Mi       799Mi       3.3Gi
Swap:            0B          0B          0B

# 关闭防火墙
sudo ufw disable

## 验证
sudo ufw status
Status: inactive

# 将桥接的 IPv4/IPv6 流量传递到 iptables 的链
```bash
## 系统启动时的模块加载配置文件, 此处新建k8s.conf,开启bridge-netfilter功能
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

## 内核配置加载,开启在ip6tables和iptables链上开启包过滤功能
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
```

2. 安装docker
```bash
# 安装curl及附属https工具
sudo apt install -y apt-transport-https ca-certificates curl

# docker gpg证书安装
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - 

# 追加docker历史源
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# 查看历史docker版本
sudo apt-cache policy docker-ce

# 安装指定版本
sudo apt install docker-ce=5:20.10.7~3-0~ubuntu-focal -y

# 修改cgroup驱动为systemd及镜像源选择aliyuncs
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://uy35zvn6.mirror.aliyuncs.com"],
  "exec-opts": ["native.cgroupdriver=systemd"]
}
EOF

# 重新加载systemd配置
sudo systemctl daemon-reload

# 启动docker
sudo systemctl restart docker

## 验证
sudo docker info | grep -i cgroup
 Cgroup Driver: systemd
 Cgroup Version: 1

## 用户追加所属组
sudo usermod -aG docker alfie
```

3. 安装kubelet, kubeadm 和 kubctl
```bash
# 添加GPG密钥以验证软件包身份和完整性
sudo curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -

# 添加 aliyun的Kubernetes apt 源 
sudo tee /etc/apt/sources.list.d/kubernetes.list <<-'EOF'
deb https://mirrors.aliyun.com/kubernetes/apt kubernetes-xenial main
EOF

# 更新源
sudo apt update

# 安装三大kube工具
# kubelet: agent进程用于与kubernetes API通信,执行任务
# kubeadm: 安装/更新kubernetes的工具
# kubectl: 与kubernetes通信的命令行工具
sudo apt-get install -y kubelet=1.22.2-00 kubeadm=1.22.2-00 kubectl=1.22.2-00 

# 对版本进行锁定不更新
sudo apt-mark hold kubelet kubeadm kubectl
```

4. Master节点进行初始化
```bash
sudo kubeadm init \
 --image-repository registry.aliyuncs.com/google_containers \
 --kubernetes-version v1.22.2 \
 --pod-network-cidr=10.244.0.0/16

[init] Using Kubernetes version: v1.22.2
[preflight] Running pre-flight checks
[preflight] Pulling images required for setting up a Kubernetes cluster
[preflight] This might take a minute or two, depending on the speed of your internet connection
[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
[certs] Using certificateDir folder "/etc/kubernetes/pki"
[certs] Generating "ca" certificate and key
[certs] Generating "apiserver" certificate and key
[certs] apiserver serving cert is signed for DNS names [kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local m1] and IPs [10.96.0.1 192.168.100.120]
[certs] Generating "apiserver-kubelet-client" certificate and key
[certs] Generating "front-proxy-ca" certificate and key
[certs] Generating "front-proxy-client" certificate and key
[certs] Generating "etcd/ca" certificate and key
[certs] Generating "etcd/server" certificate and key
[certs] etcd/server serving cert is signed for DNS names [localhost m1] and IPs [192.168.100.120 127.0.0.1 ::1]
[certs] Generating "etcd/peer" certificate and key
[certs] etcd/peer serving cert is signed for DNS names [localhost m1] and IPs [192.168.100.120 127.0.0.1 ::1]
[certs] Generating "etcd/healthcheck-client" certificate and key
[certs] Generating "apiserver-etcd-client" certificate and key
[certs] Generating "sa" key and public key
[kubeconfig] Using kubeconfig folder "/etc/kubernetes"
[kubeconfig] Writing "admin.conf" kubeconfig file
[kubeconfig] Writing "kubelet.conf" kubeconfig file
[kubeconfig] Writing "controller-manager.conf" kubeconfig file
[kubeconfig] Writing "scheduler.conf" kubeconfig file
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Starting the kubelet
[control-plane] Using manifest folder "/etc/kubernetes/manifests"
[control-plane] Creating static Pod manifest for "kube-apiserver"
[control-plane] Creating static Pod manifest for "kube-controller-manager"
[control-plane] Creating static Pod manifest for "kube-scheduler"
[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
[apiclient] All control plane components are healthy after 6.504622 seconds
[upload-config] Storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
[kubelet] Creating a ConfigMap "kubelet-config-1.22" in namespace kube-system with the configuration for the kubelets in the cluster
[upload-certs] Skipping phase. Please see --upload-certs
[mark-control-plane] Marking the node m1 as control-plane by adding the labels: [node-role.kubernetes.io/master(deprecated) node-role.kubernetes.io/control-plane node.kubernetes.io/exclude-from-external-load-balancers]
[mark-control-plane] Marking the node m1 as control-plane by adding the taints [node-role.kubernetes.io/master:NoSchedule]
[bootstrap-token] Using token:  2w2fdx.bp9ya6es2kr7349m
[bootstrap-token] Configuring bootstrap tokens, cluster-info ConfigMap, RBAC Roles
[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens to get nodes
[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
[bootstrap-token] configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
[bootstrap-token] configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
[bootstrap-token] Creating the "cluster-info" ConfigMap in the "kube-public" namespace
[kubelet-finalize] Updating "/etc/kubernetes/kubelet.conf" to point to a rotatable kubelet client certificate and key
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy

Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 192.168.100.120:6443 --token 2w2fdx.bp9ya6es2kr7349m \
	--discovery-token-ca-cert-hash sha256:bdb91eb94f2805cd66b63d837045281ca862748f4fdb58815d1c977470c47648 

# 复制登录凭据到用户目录
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# 验证
# NotReady是因为网络插件未安装
kubectl get nodes
NAME   STATUS     ROLES                  AGE   VERSION
m1     NotReady   control-plane,master   77s   v1.22.2
```

5. master节点安装CNI
```bash
# 使master节点可调度，移除taint
kubectl taint nodes --all node-role.kubernetes.io/master-

# calico离线安装[optional]
    >>  ## 下载calico 3.22.0 docker images [optional]
    >>  https://github.com/projectcalico/calico/releases/tag/v3.22.0
    >>
    >>  ## 解压 release包, impot image [optional]
    >>  tar -xvf release-v3.22.0.tgz
    >>  cd release-v3.22.0/images
    >>
    >>  docker import calico-kube-controllers.tar calico/kube-controllers:v3.22.0
    >>  docker import calico-cni.tar  calico/cni:v3.22.0
    >>  docker import calico-pod2daemon.tar  calico/pod2daemon-flexvol:v3.22.0
    >>  docker import calico-node.tar calico/node:v3.22.0

## 下载配置文件
wget https://projectcalico.docs.tigera.io/archive/v3.22/manifests/tigera-operator.yaml
wget https://projectcalico.docs.tigera.io/archive/v3.22/manifests/custom-resources.yaml

# 由于kubeadm初始化对pod网络进行了修改此处也需要进行修改
# 默认的192.168.0.0 修改为 10.244.0.0
sed  -i 's!192\.168\.0\.\0!10\.244\.0\.0!' custom-resources.yaml

# 应用配置文件
kubectl create -f tigera-operator.yaml
kubectl create -f custom-resources.yaml

# 等待docker pull images, 在此期间可以通过docker images 查看
# calico/pod2daemon-flexvol 默认会下载 v3.24.3 启动失败
# 可以通过docker 关闭相应container, 通过 docker rmi calico/pod2daemon-flexvol:v3.24.3
# 下载v3.23.5版本 docker pull calico/pod2daemon-flexvol:v3.23.5
docker images (quay和calicoi的mage为应用上述文件后获得)
REPOSITORY                                                        TAG        IMAGE ID       CREATED         SIZE
quay.io/tigera/operator                                           v1.27.16   2e97d01c05ca   9 months ago    64.5MB
calico/typha                                                      v3.23.5    9b502c848dc9   10 months ago   122MB
calico/kube-controllers                                           v3.23.5    ea5536b1fa4a   10 months ago   127MB
calico/apiserver                                                  v3.23.5    3e9e5c730042   10 months ago   84.7MB
calico/cni                                                        v3.23.5    1c979d623de9   10 months ago   254MB
calico/pod2daemon-flexvol                                         v3.23.5    5af1fa2cbe5c   10 months ago   18.8MB
calico/node                                                       v3.23.5    b6e6ee0788f2   10 months ago   207MB
registry.aliyuncs.com/google_containers/kube-apiserver            v1.22.2    e64579b7d886   23 months ago   128MB
registry.aliyuncs.com/google_containers/kube-controller-manager   v1.22.2    5425bcbd23c5   23 months ago   122MB
registry.aliyuncs.com/google_containers/kube-scheduler            v1.22.2    b51ddc1014b0   23 months ago   52.7MB
registry.aliyuncs.com/google_containers/kube-proxy                v1.22.2    873127efbc8a   23 months ago   104MB
registry.aliyuncs.com/google_containers/etcd                      3.5.0-0    004811815584   2 years ago     295MB
registry.aliyuncs.com/google_containers/coredns                   v1.8.4     8d147537fb7d   2 years ago     47.6MB
registry.aliyuncs.com/google_containers/pause                     3.5        ed210e3e4a5b   2 years ago     683kB

## 查看节点状态
kubectl get nodes 
NAME   STATUS   ROLES                  AGE   VERSION
m1     Ready    control-plane,master   39m   v1.22.2
```

6. 检查CS状态
```bash
## 查看CS状态
kubectl get cs 
Warning: v1 ComponentStatus is deprecated in v1.19+
NAME                 STATUS      MESSAGE                                                                                       ERROR
scheduler            Unhealthy   Get "http://127.0.0.1:10251/healthz": dial tcp 127.0.0.1:10251: connect: connection refused   
controller-manager   Healthy     ok                                                                                            
etcd-0               Healthy     {"health":"true","reason":""}         

## 将kube-controller-manager.yaml和kube-scheduler.yaml 的 port=0注释掉
sudo sed -i '/--port=0/ s!\(.*\)!#\1!' /etc/kubernetes/manifests/kube-controller-manager.yaml
sudo sed -i '/--port=0/ s!\(.*\)!#\1!' /etc/kubernetes/manifests/kube-scheduler.yaml 

## 重启master节点上的kubelet
sudo systemctl restart kubelet.service

## 验证scheduler状态
kubectl get cs 
Warning: v1 ComponentStatus is deprecated in v1.19+
NAME                 STATUS    MESSAGE                         ERROR
scheduler            Healthy   ok                              
controller-manager   Healthy   ok                              
etcd-0               Healthy   {"health":"true","reason":""}   
```

7. 从节点加入
```bash
# 从节点加入集群
sudo kubeadm join 192.168.100.120:6443 --token 2w2fdx.bp9ya6es2kr7349m \
	--discovery-token-ca-cert-hash sha256:bdb91eb94f2805cd66b63d837045281ca862748f4fdb58815d1c977470c47648 

# 将master的taint恢复
kubectl taint nodes m1 node-role.kubernetes.io/master=:NoSchedule
```

8. 安装自动补齐查看pod
```bash
# 自动补齐工具安装 master
source <(kubectl completion bash)
echo "source <(kubectl completion bash)" >> ~/.bashrc

kubectl get pod -A 
NAMESPACE          NAME                                      READY   STATUS    RESTARTS   AGE
calico-apiserver   calico-apiserver-5f8f4cb658-7rtnb         1/1     Running   0          30m
calico-apiserver   calico-apiserver-5f8f4cb658-f4kdh         1/1     Running   0          30m
calico-system      calico-kube-controllers-988c95d46-86kbq   1/1     Running   0          44m
calico-system      calico-node-kdtkl                         1/1     Running   0          15m
calico-system      calico-node-ltwdb                         1/1     Running   0          15m
calico-system      calico-node-vfr9k                         1/1     Running   0          44m
calico-system      calico-typha-76784cdc5b-6h5dq             1/1     Running   0          15m
calico-system      calico-typha-76784cdc5b-nccxw             1/1     Running   0          44m
kube-system        coredns-7f6cbbb7b8-9zqpb                  1/1     Running   0          62m
kube-system        coredns-7f6cbbb7b8-fk4xc                  1/1     Running   0          62m
kube-system        etcd-m1                                   1/1     Running   1          62m
kube-system        kube-apiserver-m1                         1/1     Running   1          62m
kube-system        kube-controller-manager-m1                1/1     Running   1          61m
kube-system        kube-proxy-fxrdd                          1/1     Running   0          15m
kube-system        kube-proxy-hq9nm                          1/1     Running   0          15m
kube-system        kube-proxy-p2s8s                          1/1     Running   1          62m
kube-system        kube-scheduler-m1                         1/1     Running   1          61m
tigera-operator    tigera-operator-9f6fb5887-mj77x           1/1     Running   0          46m

kubectl  get nodes 
NAME   STATUS   ROLES                  AGE   VERSION
m1     Ready    control-plane,master   63m   v1.22.2
n1     Ready    <none>                 16m   v1.22.2
n2     Ready    <none>                 16m   v1.22.2
```

### ref
[How to Run Kubernetes with Calico](https://phoenixnap.com/kb/calico-kubernetes)  
[What is difference between the options "autoclean", "autoremove" and "clean"?](https://askubuntu.com/a/3169)  
[Installing kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)  
[基于vmware16 和 ubuntu20.04, 搭建单节点 kubernetes 1.22.2](https://glory.blog.csdn.net/article/details/120606787)  
[Quickstart for Calico on Kubernetes](https://docs.tigera.io/calico/latest/getting-started/kubernetes/quickstart)  
[tigera calico 3.23](https://docs.tigera.io/archive/v3.23/about/about-calico)  
[calico github release 3.23.3](https://github.com/projectcalico/calico/releases/tag/v3.23.3)  
[kubectl cheatsheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)  

---

# 课后练习 4.2

## 启动一个 Envoy Deployment。
## 要求 Envoy 的启动配置从外部的配置文件 Mount 进 Pod。
## 进入 Pod 查看 Envoy 进程和配置。
## 更改配置的监听端口并测试访问入口的变化。
## 通过非级联删除的方法逐个删除对象。

* 创建标签为run=envoy和name=envoy副本数量为1的deployment
* 选择器选择run=envoy的容器
* 构建容器的模板,将会为容器选择envoyproxy/envoy-dev的镜像并附上标签run=envoy
* 挂载上只读的mountPath卷于/etc/envoy
* 卷的配置信息configMap中envoy-config内
```yaml
# envoy-deploy.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: envoy
  name: envoy
spec:
  replicas: 1
  selector:
    matchLabels:
      run: envoy
  template:
    metadata:
      labels:
        run: envoy
    spec:
      containers:
      - image: envoyproxy/envoy-dev
        name: envoy
        volumeMounts:
        - name: envoy-config
          mountPath: "/etc/envoy"
          readOnly: true
      volumes:
      - name: envoy-config
        configMap:
          name: envoy-config
```

* 将配置文件传入envoy-config, 吻合上述挂载要求
* 使得configMap的envoy-config变量保存envoy的配置信息
```bash
kubectl create configmap envoy-config --from-file=envoy.yaml
```

* 创建deploy
```bash
kubectl create -f envoy-deploy.yaml
```

* 访问端口
```bash
kubectl get pods -owide 
NAME                    READY   STATUS      RESTARTS   AGE   IP              NODE   NOMINATED NODE   READINESS GATES
envoy-fb5d77cc9-xcw8w   1/1     Running     0          21m   10.244.217.1    n2     <none>           <none>

curl 10.244.217.1:10000
no healthy upstream
```

* 修改configMap的envoy-config端口信息
* 查看pod配置信息也出现了同步更新,但访问依然不生效怀疑是envoy需要重新加载
* 通过replicas对envoy进行扩容,创建出来的新pod变更了访问端口
```bash
kubectl edit configmaps envoy-config
            socket_address: { address: 0.0.0.0, port_value: 10000 }
       ==>  socket_address: { address: 0.0.0.0, port_value: 12345 }

kubectl exec envoy-fb5d77cc9-xcw8w -- sed -n '/address: 0.0.0.0/p' /etc/envoy/envoy.yaml
        socket_address: { address: 0.0.0.0, port_value: 12345 }

curl 10.244.217.1:12345
curl: (7) Failed to connect to 10.244.217.1 port 12345: Connection refused

kubectl scale --replicas=2 deployment envoy

kubectl get pod -owide 
NAME                    READY   STATUS      RESTARTS   AGE     IP              NODE   NOMINATED NODE   READINESS GATES
envoy-fb5d77cc9-5hzjt   1/1     Running     0          5m11s   10.244.40.133   n1     <none>           <none>
envoy-fb5d77cc9-xcw8w   1/1     Running     0          37m     10.244.217.1    n2     <none>           <none>

curl 10.244.40.133:10000
curl: (7) Failed to connect to 10.244.40.133 port 10000: Connection refused
curl 10.244.40.133:12345
no healthy upstreama
```

* 非级联删除对象
* 删除configMap对象
* 扩容deploy replicas=3, pod由于缺少配置文件创建失败
* ---
* 删除pod对象
* 由于configMap不存在创建失败
* ---
* 删除deployment
```bash
kubectl delete configmaps envoy-config 
configmap "envoy-config" deleted

kubectl scale --replicas=3 deployment envoy 
deployment.apps/envoy scaled

kubectl get pod envoy-fb5d77cc9-7758k 
NAME                    READY   STATUS        RESTARTS   AGE
envoy-fb5d77cc9-7758k   0/1     Terminating   0          2m
Events:
  Type     Reason       Age                From               Message
  ----     ------       ----               ----               -------
  Normal   Scheduled    89s                default-scheduler  Successfully assigned default/envoy-fb5d77cc9-7758k to n2
  Warning  FailedMount  25s (x8 over 89s)  kubelet            MountVolume.SetUp failed for volume "envoy-config" : configmap "envoy-config" not found
---
kubectl delete pod envoy-fb5d77cc9-xcw8w envoy-fb5d77cc9-5hzjt 

kubectl get deploy 
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
envoy   0/2     2            0           47m
---
kubectl delete deployment envoy
deployment.apps "envoy" deleted
```