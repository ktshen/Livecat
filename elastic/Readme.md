# 啟動eserver
```
sudo gunicorn -b 120.126.16.88:17777 -k gevent -w 32 eserver:app
```

# 啟動 dbmanager
```
python3 dbmanager.py
```

# 啟動elastic docker
`sudo docker-compose -f elastic-docker.yml up`



# Elasticsearch Install and Configure

## Install

[Reference](https://www.elastic.co/guide/en/elasticsearch/reference/current/targz.html)

```
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.1.1-linux-x86_64.tar.gz &&
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.1.1-linux-x86_64.tar.gz.sha512 &&
shasum -a 512 -c elasticsearch-7.1.1-linux-x86_64.tar.gz.sha512 &&
tar -xzf elasticsearch-7.1.1-linux-x86_64.tar.gz &&
cd elasticsearch-7.1.1/
```


## Create Systemd Service and edit limit

```
# sudo vim /etc/systemd/system/elasticsearch.service
[Unit]
Description=Elasticsearch
Documentation=http://www.elastic.co
Wants=network-online.target
After=network-online.target

[Service]
RuntimeDirectory=elasticsearch
PrivateTmp=true
# Make sure that these directory is correct
Environment=ES_HOME=/home/user/elasticsearch-7.1.1
Environment=ES_PATH_CONF=/home/user/elasticsearch-7.1.1/config
Environment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin

WorkingDirectory=/home/user/elasticsearch-7.1.1

User=user
Group=user

ExecStart=/home/user/elasticsearch-7.1.1/bin/elasticsearch -p /home/user/elasticsearch-7.1.1/elasticsearch.pid --quiet

# StandardOutput is configured to redirect to journalctl since
# some error messages may be logged in standard output before
# elasticsearch logging system is initialized. Elasticsearch
# stores its logs in /var/log/elasticsearch and does not use
# journalctl by default. If you also want to enable journalctl
# logging, you can simply remove the "quiet" option from ExecStart.
StandardOutput=journal
StandardError=inherit

# Specifies the maximum file descriptor number that can be opened by this process
LimitNOFILE=65535

# Specifies the maximum number of processes
LimitNPROC=4096

# Specifies the maximum size of virtual memory
LimitAS=infinity

# Specifies the maximum file size
LimitFSIZE=infinity

# Disable timeout logic and wait until process is stopped
TimeoutStopSec=0

# SIGTERM signal is used to stop the Java process
KillSignal=SIGTERM

# Send the signal only to the JVM rather than its control group
KillMode=process

# Lock memory
LimitMEMLOCK=infinity

# Java process is never killed
SendSIGKILL=no

# When a JVM receives a SIGTERM signal it exits with code 143
SuccessExitStatus=143

[Install]
WantedBy=multi-user.target
```

```
# sudo vim /etc/security/limits.conf
# Reboot to make the setting effective
user soft memlock unlimited   # user is the user name
user hard memlock unlimited
```

## Configure

```
# vim config/elasticsearch.yml
cluster.name: livenet-cluster
# This field should be unique
node.name: 120-126-16-88-node
# Make sure that the user gets the access to the data directory
path.data: /var/lib/elasticsearch
path.logs: /var/log/elasticsearch
bootstrap.memory_lock: true
network.host: 0.0.0.0
# host external ip
network.publish_host: 120.126.16.88
http.port: 9200
# Add every node's ip and name to the following fields
discovery.seed_hosts: ["120.126.16.88"]
cluster.initial_master_nodes: ["120-126-16-88-node"]
```

```
# vim config/jvm.options
# set the space used by the heap in memory
-Xms2g
-Xmx2g
```

## Command


```
# Enable Service
sudo /bin/systemctl daemon-reload
sudo /bin/systemctl enable elasticsearch.service

# Service command
sudo systemctl start elasticsearch.service
sudo systemctl status elasticsearch.service
sudo systemctl stop elasticsearch.service

# logging
sudo journalctl -f
sudo journalctl --unit elasticsearch
```


## Check cluster health

```
curl -X GET "localhost:9200/_cluster/health?pretty"
```

-   Make sure that the status is green

## Create new index

```
curl -X PUT "localhost:9200/livestreams" -H 'Content-Type: application/json' -d' { "settings" : { "index" : { "number_of_shards" : 3, "number_of_replicas" : 1 } } , "mappings" : { "properties" : { "timestamp" : { "type" : "date" } , "published" : { "type" : "date" } , "click_through" : { "type" : "integer" } , "viewers" : { "type" : "integer" } , "viewcount" : { "type" : "integer" } , "popular_rate": { "type" : "integer" } } } } '
```

