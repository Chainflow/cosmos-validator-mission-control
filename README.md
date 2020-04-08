# Validator monitoring and alerting tool

## Prerequisites
- **Go 13.x+**

**Install Grafana on ununtu**

Download the latest .tar.gz file and extract it by using the following commands

```sh
$ cd $HOME
$ wget https://dl.grafana.com/oss/release/grafana-6.7.2.linux-amd64.tar.gz
$ tar -zxvf grafana-6.7.2.linux-amd64.tar.gz
```

Start the grafana server
```sh
$ cd grafana-6.7.2/bin/
$ ./grafana-server
```

**Install InfluxDB**

Download the latest .tar.gz file and extract it by using the following commands

```sh
$ cd $HOME
$ wget https://dl.influxdata.com/influxdb/releases/influxdb-1.7.10_linux_amd64.tar.gz
$ tar xvfz influxdb-1.7.10_linux_amd64.tar.gz
```

Start influxDB

```sh
$ cd $HOME and run the below command to start the server
$ ./influxdb-1.7.10-1/usr/bin/influxd
```

**Note :** If you want to give custom configuration then you can edit the `influxdb.conf` at `/influxdb-1.7.10-1/etc/influxdb` and do not forget to restart the server after the changes.


**Telegraf Installation**

Download the latest .tar.gz file and extract it by using the following commands
```sh
$ cd $HOME
$ wget https://dl.influxdata.com/telegraf/releases/telegraf-1.14.0_linux_amd64.tar.gz
tar xf telegraf-1.14.0_linux_amd64.tar.gz
```

Start telegraph
```sh
$ cd telegraf/usr/bin/
$ ./telegraf --config ../../etc/telegraf/telegraf.conf
```

## Get the code
```bash
$ git clone git@github.com:chris-remus/chainflow-vitwit.git
$ cd chainflow-vitwit
$ git fetch && git checkout releases/alpha
$ cp example.config.toml config.toml
```

`config.toml` has following configurations
- *tg_chat_id*

    Telegram chat ID to receive telegram alerts
- *tg_bot_token*

    Telegram bot token. The bot should be added to the chat and it should have send msessage permission

- *email_address*

    Email address to receive email notifications

- *sendgrid_token*

    Sendgrid mail service api token.
- *missed_blocks_threshold*

    Configure the threshold to receive  **Missed Block Alerting**
- *block_diff_threshold*

    An integer value to receive **Block difference alerts**

- *alert_time1* and *alert_time2*

    These are for regular status updates. To receive validator status daily (twice), configure these parameters in the form of "02:25PM". The time here refers to UTC time.

After populating config.toml, build and run the monitoring binary

```bash
$ go build -o chain-monit && ./chain-monit
```

Run using docker
```bash
$ docker build -t cfv .
$ docker run -d --name chain-monit cfv
```
