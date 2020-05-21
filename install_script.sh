#!/bin/bash

cd $HOME

export INFLUX_DAEMON=influxd
export TELEGRAF=telegraf

teleFalg="$1"
teleFlagValue="--remote-hosted"

echo "----------- Installing grafana -----------"

sudo apt-get install -y adduser libfontconfig1

wget https://dl.grafana.com/oss/release/grafana_6.7.2_amd64.deb

sudo dpkg -i grafana_6.7.2_amd64.deb

echo "------ Starting grafana server using systemd --------"

sudo systemctl daemon-reload

sudo systemctl start grafana-server

sudo systemctl status grafana-server

cd $HOME

echo "----------- Installing Influx -----------"

wget https://dl.influxdata.com/influxdb/releases/influxdb-1.7.10_linux_amd64.tar.gz

tar xvfz influxdb-1.7.10_linux_amd64.tar.gz

echo "-----------Fetching influxd path--------"

INFLUXD_PATH=$(which $INFLUX_DAEMON)

echo "---------Creating influxd system file---------"

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

echo "------- Running influxd by using systemd ------"

sudo mv $INFLUX_DAEMON.service /lib/systemd/system/$INFLUX_DAEMON.service
sudo -S systemctl daemon-reload
sudo -S systemctl start $INFLUX_DAEMON

cd $HOME

if [ "$teleFalg" != "$teleFlagValue" ];
then 
	echo "----------- Installing telegraf -----------------"
	wget https://dl.influxdata.com/telegraf/releases/telegraf-1.14.0_linux_amd64.tar.gz

	tar xf telegraf-1.14.0_linux_amd64.tar.gz
	
	echo "-----------Fetching telegraf path--------"

	TELEGRAF_PATH=telegraf/usr/bin/telegraf

	echo "---------Creating telegraf system file---------"

	echo "[Unit]
	Description=${TELEGRAF} daemon
	After=network-online.target
	[Service]
	User=${USER}
	ExecStart=${HOME}/telegraf/usr/bin/telegraf --config telegraf/etc/telegraf/telegraf.conf
	Restart=always
	RestartSec=3
	LimitNOFILE=4096
	[Install]
	WantedBy=multi-user.target
	" >$TELEGRAF.service

	echo "---------- Running telegraf by using systemd ---------"

	sudo mv $TELEGRAF.service /lib/systemd/system/$TELEGRAF.service
	sudo -S systemctl daemon-reload
	sudo -S systemctl start $TELEGRAF

else
	echo "--remote-hosted enabled, so not downloading the telegraf"
fi

echo "--------- Cloning cosmos-validator-mission-control -----------"

git clone git@github.com:Chainflow/cosmos-validator-mission-control.git

cd cosmos-validator-mission-control

cp example.config.toml config.toml

echo "-------- Creating database in influx ----------"

influx

CREATE DATABASE vcf

exit

echo "------ Building and running the code --------"

go build && ./cosmos-validator-mission-control &

docker build -t cfv .

docker run -d --name cosmos-validator-mission-control cfv