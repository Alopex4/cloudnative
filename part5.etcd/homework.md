# 课后练习 5.1

## 按照课上讲解的方法在本地构建一个单节点的基于 HTTPS 的 etcd 集群
```bash
# 安装etcd v3.4.17
ETCD_VER=v3.4.17
DOWNLOAD_URL=https://github.com/etcd-io/etcd/releases/download
sudo curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
sudo tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
sudo rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

# 启动etcd
# 安装证书工具包: CFSSL: Cloudflare's PKI and TLS toolkit
sudo apt install golang-cfssl

# 下载etcd-io工具
mkdir /tmp/etcd-io
git clone https://github.com/etcd-io/etcd.git
# 保留127.0.0.1和localhost
vim /tmp/etcd-io/etcd/hack/tls-setup/config/req-csr.json
{
  "CN": "example.com",
  "hosts": [
    "127.0.0.1",
    "localhost"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "ST": "California",
      "L": "San Francisco",
      "O": "Example Organization",
      "OU": "IT Department"
    }
  ]
}

# 全局化变量
export infra0=127.0.0.1
export infra1=127.0.0.1
export infra2=127.0.0.1
# 验证
env | grep infra 
infra2=127.0.0.1
infra1=127.0.0.1
infra0=127.0.0.1

# 编译创建证书
cd /tmp/etcd-io/etcd/hack/tls-setup
make 
#验证
ls certs
127.0.0.1.csr  127.0.0.1-key.pem  127.0.0.1.pem  ca.csr  ca-key.pem  ca.pem  peer-127.0.0.1.csr  peer-127.0.0.1-key.pem  peer-127.0.0.1.pem

# 确认端口未被占用
ss -tulnp  | grep [3-5]380 && echo OCCUPY || echo non-OCCUPY
non-OCCUPY

# 启用etcd
mkdir /tmp/etcd-download-test/log
# 在本地开启etcd,通过不同的端口区分服务
# etcd默认为前台执行,通过nohub可将其放在后台运行不再与终端关联
# --data-dir etcd数据存放目录
# --listen-peer-urls peer节点通信地址
# --initial-advertise-peer-urls 初始化peer通信地址
# --listen-client-urls 监听客户与etcd通信地址 (服务地址/端口)
# --advertise-client-urls  客户与etcd通信地址
# --initial-cluster-token etcd集群token
# --initial-cluster 初始化集群信息
# --initial-cluster-state 初始化集群状态
# --client-cert-auth --trusted-ca-file 信任CA验证文件位置
# --cert-file 证书文件位置
# --key-file 密钥文件位置
# --peer-client-cert-auth --peer-trusted-ca-file peer节点信任CA
# --peer-cert-file peer证书文件位置
# --peer-key-file peer密钥文件位置 
# 2>&1 > /tmp/etcd-download-test/log/infra0.log & >> fd=2(stderr) 重定向到 fd=1(stdout), 将fd=1&fd=2 记录到日志/tmp/etcd-download-test/log/infra0.log
nohup /tmp/etcd-download-test/etcd --name infra0 \
--data-dir=/tmp/etcd/infra0 \
--listen-peer-urls https://127.0.0.1:3380 \
--initial-advertise-peer-urls https://127.0.0.1:3380 \
--listen-client-urls https://127.0.0.1:3379 \
--advertise-client-urls https://127.0.0.1:3379 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster infra0=https://127.0.0.1:3380,infra1=https://127.0.0.1:4380,infra2=https://127.0.0.1:5380 \
--initial-cluster-state new \
--client-cert-auth --trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--peer-cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--peer-key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem 2>&1 > /tmp/etcd-download-test/log/infra0.log  &

nohup /tmp/etcd-download-test/etcd --name infra1 \
--data-dir=/tmp/etcd/infra1 \
--listen-peer-urls https://127.0.0.1:4380 \
--initial-advertise-peer-urls https://127.0.0.1:4380 \
--listen-client-urls https://127.0.0.1:4379 \
--advertise-client-urls https://127.0.0.1:4379 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster infra0=https://127.0.0.1:3380,infra1=https://127.0.0.1:4380,infra2=https://127.0.0.1:5380 \
--initial-cluster-state new \
--client-cert-auth --trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--peer-cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--peer-key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem 2>&1 > /tmp/etcd-download-test/log/infra1.log &

nohup /tmp/etcd-download-test/etcd --name infra2 \
--data-dir=/tmp/etcd/infra2 \
--listen-peer-urls https://127.0.0.1:5380 \
--initial-advertise-peer-urls https://127.0.0.1:5380 \
--listen-client-urls https://127.0.0.1:5379 \
--advertise-client-urls https://127.0.0.1:5379 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster infra0=https://127.0.0.1:3380,infra1=https://127.0.0.1:4380,infra2=https://127.0.0.1:5380 \
--initial-cluster-state new \
--client-cert-auth --trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--peer-cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--peer-key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem 2>&1 > /tmp/etcd-download-test/log/infra2.log &

# 验证端口
ss -tulnp  | grep [3-5]380 | column  -t 
tcp  LISTEN  0  4096  127.0.0.1:4380  0.0.0.0:*  users:(("etcd",pid=1875486,fd=5))
tcp  LISTEN  0  4096  127.0.0.1:5380  0.0.0.0:*  users:(("etcd",pid=1876920,fd=5))
tcp  LISTEN  0  4096  127.0.0.1:3380  0.0.0.0:*  users:(("etcd",pid=1874234,fd=5))

# 验证成员在线
# 这里此前使用了证书验证, 因此客户端访问也需要添加证书参数
/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem --write-out=table member list 
+------------------+---------+--------+------------------------+------------------------+------------+
|        ID        | STATUS  |  NAME  |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
+------------------+---------+--------+------------------------+------------------------+------------+
| 1701f7e3861531d4 | started | infra0 | https://127.0.0.1:3380 | https://127.0.0.1:3379 |      false |
| 6a58b5afdcebd95d | started | infra1 | https://127.0.0.1:4380 | https://127.0.0.1:4379 |      false |
| 84a1a2f39cda4029 | started | infra2 | https://127.0.0.1:5380 | https://127.0.0.1:5379 |      false |
+------------------+---------+--------+------------------------+------------------------+------------+
```
## 写一条数据 (多条)
```bash
/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem put aliyun acp

/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem put google gcp

/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem put amazon aws
```

## 查看数据细节
```bash
/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem get google
google
gcp
```

## 删除数据
```bash
# 备份数据
/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem snapshot save /tmp/etcd-download-test/snapshot.db
{"level":"info","ts":1693663371.0530963,"caller":"snapshot/v3_snapshot.go:119","msg":"created temporary db file","path":"/tmp/etcd-download-test/snapshot.db.part"}
{"level":"info","ts":"2023-09-02T14:02:51.063Z","caller":"clientv3/maintenance.go:200","msg":"opened snapshot stream; downloading"}
{"level":"info","ts":1693663371.0637405,"caller":"snapshot/v3_snapshot.go:127","msg":"fetching snapshot","endpoint":"https://127.0.0.1:3379"}
{"level":"info","ts":"2023-09-02T14:02:51.067Z","caller":"clientv3/maintenance.go:208","msg":"completed snapshot read; closing"}
{"level":"info","ts":1693663371.0689585,"caller":"snapshot/v3_snapshot.go:142","msg":"fetched snapshot","endpoint":"https://127.0.0.1:3379","size":"20 kB","took":0.014806864}
{"level":"info","ts":1693663371.069047,"caller":"snapshot/v3_snapshot.go:152","msg":"saved","path":"/tmp/etcd-download-test/snapshot.db"}
Snapshot saved at /tmp/etcd-download-test/snapshot.db

# 删除数据
ls /tmp/etcd
infra0  infra1  infra2

rm -rf /tmp/etcd/*

# 访问失效(需要数十秒)
/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem get google
{"level":"warn","ts":"2023-09-02T14:06:16.497Z","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-4cf1b078-bfdc-4fb0-8219-239cc252e7ae/127.0.0.1:3379","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = latest balancer error: all SubConns are in TransientFailure, latest connection error: connection error: desc = \"transport: Error while dialing dial tcp 127.0.0.1:3379: connect: connection refused\""}
Error: context deadline exceeded

# etcd已自动退出
ss -tulnp  | grep [3-5]380 && echo OCCUPY || echo non-OCCUPY
non-OCCUPY

# 通过快照恢复数据
export ETCDCTL_API=3
/tmp/etcd-download-test/etcdctl snapshot restore /tmp/etcd-download-test/snapshot.db \
  --name infra0 \
  --data-dir=/tmp/etcd/infra0 \
  --initial-cluster infra0=https://127.0.0.1:3380,infra1=https://127.0.0.1:4380,infra2=https://127.0.0.1:5380 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-advertise-peer-urls https://127.0.0.1:3380

/tmp/etcd-download-test/etcdctl snapshot restore /tmp/etcd-download-test/snapshot.db \
    --name infra1 \
    --data-dir=/tmp/etcd/infra1 \
    --initial-cluster infra0=https://127.0.0.1:3380,infra1=https://127.0.0.1:4380,infra2=https://127.0.0.1:5380 \
    --initial-cluster-token etcd-cluster-1 \
    --initial-advertise-peer-urls https://127.0.0.1:4380

/tmp/etcd-download-test/etcdctl snapshot restore /tmp/etcd-download-test/snapshot.db \
  --name infra2 \
  --data-dir=/tmp/etcd/infra2 \
  --initial-cluster infra0=https://127.0.0.1:3380,infra1=https://127.0.0.1:4380,infra2=https://127.0.0.1:5380 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-advertise-peer-urls https://127.0.0.1:5380

{"level":"info","ts":1693663790.85998,"caller":"snapshot/v3_snapshot.go:296","msg":"restoring snapshot","path":"/tmp/etcd-download-test/snapshot.db","wal-dir":"/tmp/etcd/infra2/member/wal","data-dir":"/tmp/etcd/infra2","snap-dir":"/tmp/etcd/infra2/member/snap"}
{"level":"info","ts":1693663790.864269,"caller":"membership/cluster.go:392","msg":"added member","cluster-id":"c36a1e619c38211b","local-member-id":"0","added-peer-id":"1701f7e3861531d4","added-peer-peer-urls":["https://127.0.0.1:3380"]}
{"level":"info","ts":1693663790.8643684,"caller":"membership/cluster.go:392","msg":"added member","cluster-id":"c36a1e619c38211b","local-member-id":"0","added-peer-id":"6a58b5afdcebd95d","added-peer-peer-urls":["https://127.0.0.1:4380"]}
{"level":"info","ts":1693663790.8644247,"caller":"membership/cluster.go:392","msg":"added member","cluster-id":"c36a1e619c38211b","local-member-id":"0","added-peer-id":"84a1a2f39cda4029","added-peer-peer-urls":["https://127.0.0.1:5380"]}
{"level":"info","ts":1693663790.867847,"caller":"snapshot/v3_snapshot.go:309","msg":"restored snapshot","path":"/tmp/etcd-download-test/snapshot.db","wal-dir":"/tmp/etcd/infra2/member/wal","data-dir":"/tmp/etcd/infra2","snap-dir":"/tmp/etcd/infra2/member/snap"}

# 重启etcd
# 重启已有服务,不在需要init为prefix的参数
nohup /tmp/etcd-download-test/etcd --name infra0 \
--data-dir=/tmp/etcd/infra0 \
--listen-peer-urls https://127.0.0.1:3380 \
--listen-client-urls https://127.0.0.1:3379 \
--advertise-client-urls https://127.0.0.1:3379 \
--client-cert-auth --trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--peer-cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--peer-key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem 2>&1 > /tmp/etcd-download-test/log/infra0.log  &

nohup /tmp/etcd-download-test/etcd --name infra1 \
--data-dir=/tmp/etcd/infra1 \
--listen-peer-urls https://127.0.0.1:4380 \
--listen-client-urls https://127.0.0.1:4379 \
--advertise-client-urls https://127.0.0.1:4379 \
--client-cert-auth --trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--peer-cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--peer-key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem 2>&1 > /tmp/etcd-download-test/log/infra1.log &

nohup /tmp/etcd-download-test/etcd --name infra2 \
--data-dir=/tmp/etcd/infra2 \
--listen-peer-urls https://127.0.0.1:5380 \
--listen-client-urls https://127.0.0.1:5379 \
--advertise-client-urls https://127.0.0.1:5379 \
--client-cert-auth --trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem \
--peer-cert-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem \
--peer-key-file=/tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem 2>&1 > /tmp/etcd-download-test/log/infra2.log &

# 再次查询数据,确认是否恢复
/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem get google
google
gcp

/tmp/etcd-download-test/etcdctl --endpoints https://127.0.0.1:3379 --cert /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1.pem --key /tmp/etcd-io/etcd/hack/tls-setup/certs/127.0.0.1-key.pem --cacert /tmp/etcd-io/etcd/hack/tls-setup/certs/ca.pem get google -wjson  | jq 
{
  "header": {
    "cluster_id": 14081100589508862000,
    "member_id": 1657878694428226000,
    "revision": 4,
    "raft_term": 8
  },
  "kvs": [
    {
      "key": "Z29vZ2xl",
      "create_revision": 3,
      "mod_revision": 3,
      "version": 1,
      "value": "Z2Nw"
    }
  ],
  "count": 1
}
```

# 容量特点, 单个对象不超过1.5M, 默认容量为2G, 建议不超过8G
## 容量不足处理
```bash
# 限制etcd容量大小为16MB
/tmp/etcd-download-test/etcd --listen-client-urls 'http://localhost:12379' \
--advertise-client-urls 'http://localhost:12379' \
--listen-peer-urls 'http://localhost:12380' \
--initial-advertise-peer-urls 'http://localhost:12380' \
--initial-cluster 'default=http://localhost:12380' \
--quota-backend-bytes=$((16*1024*1024)) # 16MB 

# 查看当前etcd member
/tmp/etcd-download-test/etcdctl member list --write-out=table --endpoints=localhost:12379
+------------------+---------+---------+------------------------+------------------------+------------+
|        ID        | STATUS  |  NAME   |       PEER ADDRS       |      CLIENT ADDRS      | IS LEARNER |
+------------------+---------+---------+------------------------+------------------------+------------+
| c9ac9fc89eae9cf7 | started | default | http://localhost:12380 | http://localhost:12379 |      false |
+------------------+---------+---------+------------------------+------------------------+------------+

# 循环数据写入
# 显示` mvcc: database space exceeded`
while [ 1 ]; do dd if=/dev/urandom bs=1024 count=1024 | ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl --endpoints=localhost:12379 put key || break; done
...
OK
1024+0 records in
1024+0 records out
1048576 bytes (1.0 MB, 1.0 MiB) copied, 0.0298789 s, 35.1 MB/s
{"level":"warn","ts":"2023-09-02T14:41:56.220Z","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-5a837d4b-938e-469c-bf0d-9eb4466cb47b/localhost:12379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
Error: etcdserver: mvcc: database space exceeded

# 查看告警
/tmp/etcd-download-test/etcdctl alarm list  --endpoints=localhost:12379
memberID:14532165781622267127 alarm:NOSPACE 

# 查看日志
2023-09-02 14:41:56.220448 W | etcdserver: alarm NOSPACE raised by peer c9ac9fc89eae9cf7

# endpoint状态查看
/tmp/etcd-download-test/etcdctl endpoint status --endpoints=localhost:12379 --write-out=table 
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+
|    ENDPOINT     |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX |             ERRORS             |
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+
| localhost:12379 | c9ac9fc89eae9cf7 |  3.4.17 |   17 MB |      true |      false |        10 |         33 |                 33 |  memberID:14532165781622267127 |
|                 |                  |         |         |           |            |           |            |                    |                 alarm:NOSPACE  |
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+

# 清理空间
tmp/etcd-download-test/etcdctl defrag --endpoints=localhost:12379
Finished defragmenting etcd member[localhost:12379]

# 再次查看endpoint状态
# DBSIZE 已下降
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+
|    ENDPOINT     |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX |             ERRORS             |
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+
| localhost:12379 | c9ac9fc89eae9cf7 |  3.4.17 |   12 MB |      true |      false |        10 |         33 |                 33 |  memberID:14532165781622267127 |
|                 |                  |         |         |           |            |           |            |                    |                 alarm:NOSPACE  |
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+

# 删除alarm
/tmp/etcd-download-test/etcdctl alarm disarm --endpoints=localhost:12379 
memberID:14532165781622267127 alarm:NOSPACE 

# 告警清理成功
/tmp/etcd-download-test/etcdctl endpoint status --endpoints=localhost:12379 --write-out=table 
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
|    ENDPOINT     |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX | ERRORS |
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+
| localhost:12379 | c9ac9fc89eae9cf7 |  3.4.17 |   12 MB |      true |      false |        10 |         35 |                 35 |        |
+-----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------+

# 继续写入正常
/tmp/etcd-download-test/etcdctl put a b --endpoints=localhost:12379 
OK

# 压缩至版本3
/tmp/etcd-download-test/etcdctl compact 3
```