
# Initial example-complete repository

```
TODO_CHANGE_DB - db username
TODO_WEBAPP_SERVICE - the webapp
TODO_CHANGE_DOMAIN - the domain name
TODO_CHANGE_THIS - either gmail user or password
example-complete - the project name (also the $GOPATH/example-complete directory)
example-cron - the background service
```

## Server Setup
### install Ubuntu Server 17.04
```
sudo apt-get install python curl htop silversearcher-ag vnstat zsh most inxi apt-file  golang-1.8 ruby git ruby-dev zlib1g-dev ffmpeg upx pv psmisc expect time byobu libpq-dev postgresql-server-dev-all 
sudo ln -s /usr/lib/go/bin/* /usr/bin

sudo systemctl enable vnstat
sudo systemctl start vnstat
```

#### setup postgresql
```
sudo apt-get install postgresql
sudo systemctl enable postgresql
sudo sed -i 's|local   all             all                                     peer|local all all trust|g' /etc/postgresql/9.6/main/pg_hba.conf
sudo sed -i 's|host    all             all             127.0.0.1/32            md5|host all all 127.0.0.1/32 trust|g' /etc/postgresql/9.6/main/pg_hba.conf
sudo sed -i 's|host    all             all             ::1/128                 md5|host all all ::1/128 trust|g' /etc/postgresql/9.6/main/pg_hba.conf

sudo systemctl stop postgresql
sudo mv /var/lib/postgresql /home/
sudo ln -s /home/postgresql /var/lib/postgresql

sudo systemctl start postgresql # psql
sudo su - postgres <<EOF
createuser baik
createdb baik
psql -c 'GRANT ALL PRIVILEGES ON DATABASE TODO_CHANGE_DB TO TODO_CHANGE_DB;'
EOF
```

#### setup scylladb
```
echo '# http://www.scylladb.com/download/ubuntu-16-04/?mkt_tok=eyJpIjoiTldOallqTXhNVFppWVRjdyIsInQiOiIzR0huYjFmWFZhbG9Cd0VWMnp3THc0Qzl2Zk9VM0xLZ1NaMWFqcVhDb3IxU0Qwb0FiYll3VUlLcmhKKzMzek13RG5kXC8wajltNE94VE1acEZyXC9kWmFWcnJuek1EZ0J4K2NiYkswU1A0ZVM5RUZhK1wvYU9vNmZ0WEF5aWFoSzRPSyJ9
deb http://archive.ubuntu.com/ubuntu/ xenial main restricted universe multiverse
' | sudo tee /etc/apt/sources.list.d/xenial.list

sudo apt-get install software-properties-common python-software-properties
sudo add-apt-repository -y ppa:openjdk-r/ppa
sudo wget -O /etc/apt/sources.list.d/scylla.list http://downloads.scylladb.com/deb/ubuntu/scylla-1.7-xenial.list
sudo apt-add-repository ppa:arnaud-hartmann/glances-stable
sudo apt-get update

sudo apt-get install -y --allow-unauthenticated openjdk-8-jre-headless glances scylla-server scylla-jmx scylla-tools
#sudo update-java-alternatives -s java-1.8.0-openjdk-amd64
sudo systemctl enable scylla-server

sudo scylla_setup
# NO!!! It is recommended to use RAID0 and XFS for Scylla data. If you select yes, you will be prompt to choose which unmounted disks to use for Scylla data. Selected disks will be formatted in the process.
# NO!!! Do you want to setup coredump? Answer yes to enable core dumps; this allows to do post-mortem analysis of Scylla state after a crash. Answer no to do nothing.
# Generating evaluation file sized 10GB... timed out before we could write the entire file. Will continue but accuracy may suffer.
# 2.92871GB written in 144 seconds
# Refining search for maximum. So far, 3968 IOPS 
# Maximum throughput: 3968 IOPS
# NO!!! Do you want to install node exporter and export Prometheus data from the node?
sudo sed -i 's|/usr/bin/scylla $SCYLLA_ARGS|/usr/bin/scylla -m 1G -c 1 --developer-mode 1 $SCYLLA_ARGS|g' /lib/systemd/system/scylla-server.service

echo '
experimental: true
' | sudo tee -a /etc/scylla/scylla.yaml
sudo mv /var/lib/scylla /home/
sudo ln -s /home/scylla /var/lib/scylla

sudo reboot

sudo systemctl start scylla-server # cqlsh 127.0.0.1
nodetool status
```

# proxy: caddy
```
curl https://getcaddy.com | bash
sudo ln -s /usr/local/bin/caddy /usr/bin/caddy
sudo setcap 'cap_net_bind_service=+ep' `which caddy`
echo "
$(whoami) soft nofile 1048576 # default: 1024
$(whoami) hard nofile 2097152
$(whoami) soft noproc 262144  # default 128039
$(whoami) hard noproc 524288 
postgres soft nofile 1048576 # default: 1024
postgres hard nofile 2097152
postgres soft noproc 262144  # default 128039
postgres hard noproc 524288 
" | sudo tee -a /etc/security/limits.conf
cat /etc/security/limits.d/scylla.conf | sudo tee -a /etc/security/limits.conf
sudo mv /etc/security/limits.d/scylla.conf /etc/security/limits.d/scylla.conf.old
```


# tuning
```
echo '
net.ipv4.tcp_max_syn_backlog = 8192 # default: 4096
net.core.somaxconn = 8192 # default: 4096
net.ipv4.tcp_tw_reuse = 1 # default: 0
net.ipv4.tcp_tw_recycle = 1 # default: 0
' | sudo tee /etc/sysctl.d/99-kyztune.conf
sudo sysctl -p --system 
```