#!/bin/bash

cd $HOME

export INFLUX_DAEMON=influxd

teleFalg="$1"
teleFlagValue="--remote-hosted"

wget https://dl.grafana.com/oss/release/grafana-6.7.2.linux-amd64.tar.gz

tar -zxvf grafana-6.7.2.linux-amd64.tar.gz

cd grafana-6.7.2/bin/

./grafana-server &

cd $HOME

wget https://dl.influxdata.com/influxdb/releases/influxdb-1.7.10_linux_amd64.tar.gz

tar xvfz influxdb-1.7.10_linux_amd64.tar.gz

echo "-----------Fetching influxd path--------"
INFLUXD_PATH=$(which $INFLUX_DAEMON)

echo "---------Creating system file---------"

echo "[Unit]
Description=${INFLUX_DAEMON} daemon
After=network-online.target
[Service]
User=${USER}
ExecStart=${INFLUXD_PATH}
Restart=always
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target
" >$INFLUX_DAEMON.service

sudo mv $INFLUX_DAEMON.service /lib/systemd/system/$INFLUX_DAEMON.service
sudo -S systemctl daemon-reload
sudo -S systemctl start $INFLUX_DAEMON

cd $HOME

if [ "$teleFalg" != "$teleFlagValue" ];
then 
	wget https://dl.influxdata.com/telegraf/releases/telegraf-1.14.0_linux_amd64.tar.gz

	tar xf telegraf-1.14.0_linux_amd64.tar.gz

	cd telegraf/usr/bin/
 
	./telegraf --config ../../etc/telegraf/telegraf.conf &
else
	echo "--remote-hosted enabled, so not downloading the telegraf"
fi

git clone git@github.com:Chainflow/cosmos-validator-mission-control.git

cd cosmos-validator-mission-control

cp example.config.toml config.toml

influx

CREATE DATABASE vcf

exit

go build && ./cosmos-validator-mission-control &

docker build -t cfv .

docker run -d --name cosmos-validator-mission-control cfv