
# Uninstall nginx if installed

sudo apt purge -y nginx nginx-common nginx-core
sudo apt autoremove


# Misc

sudo apt install -y net-tools byobu


# Caddy

sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo tee /etc/apt/trusted.gpg.d/caddy-stable.asc
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy


# Clickhouse

sudo apt-get install -y apt-transport-https ca-certificates dirmngr
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv E0C56BD4

echo "deb https://repo.clickhouse.tech/deb/stable/ main/" | sudo tee \
    /etc/apt/sources.list.d/clickhouse.list
sudo apt-get -y update

sudo apt-get install -y clickhouse-server clickhouse-client

# TODO: edit /etc/clickhouse-server/config.xml set listen only to 127.0.0.1

sudo systemctl enable clickhouse-server
sudo systemctl start clickhouse-server

journalctl -r -u clickhouse-server

clickhouse-client # connect


# Tarantool

# old way for server

#sudo systemctl disable tarantool
#sudo systemctl stop tarantool
#sudo apt purge tarantool
#
#curl -L https://tarantool.io/release/2.8/installer.sh | bash
#sudo apt-get -y install tarantool
#
#sudo systemctl enable tarantool
#sudo systemctl start tarantool
#
#journalctl -r -u tarantool
#
#tarantoolctl connect 3301 # connect

# manual for local

apt-get -y install sudo
sudo apt-get -y install gnupg2 curl
curl https://download.tarantool.org/tarantool/release/2.8/gpgkey | sudo apt-key add -
sudo apt-get -y install lsb-release
release=`lsb_release -c -s` 
sudo apt-get -y install apt-transport-https 
sudo rm -f /etc/apt/sources.list.d/*tarantool*.list
echo "deb https://download.tarantool.org/tarantool/release/2.8/ubuntu/ ${release} main" | sudo tee /etc/apt/sources.list.d/tarantool_2_8.list
echo "deb-src https://download.tarantool.org/tarantool/release/2.8/ubuntu/ ${release} main" | sudo tee -a /etc/apt/sources.list.d/tarantool_2_8.list
sudo apt-get -y update
sudo apt-get -y install tarantool 

# multi-insance for server

wget http://ftp.au.debian.org/debian/pool/main/n/netselect/netselect_0.3.ds1-26_amd64.deb
sudo dpkg -i netselect_0.3.ds1-26_amd64.deb
sudo netselect -s 20 -t 40 $(wget -qO - mirrors.ubuntu.com/mirrors.txt)
sudo sed -i 's/archive.ubuntu.com/mirror.biznetgio.com/' /etc/apt/sources.list

sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get install ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "
deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable
" | sudo tee /etc/apt/sources.list.d/docker.list 
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io

sudo usermod -aG docker $USER
newgrp docker 
sudo systemctl enable docker.service
sudo systemctl enable containerd.service

echo 'version: "3.3"

services:

  tt1:
    container_name: tt1
    hostname: tt1
    image: tarantool/tarantool:latest # 2.8.3 
    # x.x.0 = alpha, x.x.1 = beta, x.x.2+ = stable, latest not always stable
    environment:
      TARANTOOL_USER_NAME: "myusername" 
      TARANTOOL_USER_PASSWORD: "mysecretpassword"
    volumes:
      - ./tarantool-data:/var/lib/tarantool
    ports:
      - "3301:3301"
    restart: always
' > docker-compose.yaml

sudo apt install docker-compose

docker-compose up -d

# connect manually:

tarantoolctl connect myusername:mysecretpassword@localhost:3301

docker exec -t -i tt1 console

