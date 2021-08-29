
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

sudo systemctl disable tarantool
sudo systemctl stop tarantool
sudo apt purge tarantool

curl -L https://tarantool.io/release/2.7/installer.sh | bash
sudo apt-get -y install tarantool

sudo systemctl enable tarantool
sudo systemctl start tarantool

journalctl -r -u tarantool

tarantoolctl connect 3301 # connect
