## Validator monitoring and alerting tool

### Getting Started

- Grafana Installation on ununtu

```sh
- **Download the latest .tar.gz file and extract it by using the following commands**

    cd $HOME
    wget https://dl.grafana.com/oss/release/grafana-6.7.2.linux-amd64.tar.gz
    tar -zxvf grafana-6.7.2.linux-amd64.tar.gz

- **Start the grafana server**

    cd grafana-6.7.2/bin/
    ./grafana-server

```

- InfluxDB Installation

```sh
- **Download the latest .tar.gz file and extract it by using the following commands**

    cd $HOME
    wget https://dl.influxdata.com/influxdb/releases/influxdb-1.7.10_linux_amd64.tar.gz
    tar xvfz influxdb-1.7.10_linux_amd64.tar.gz

- **Start influxDB**

    cd $HOME and run the below command to start the server
    ./influxdb-1.7.10-1/usr/bin/influxd

    **Note : If you want to give your own configuration then you can edit the influxdb.conf file which is in the path /influxdb-1.7.10-1/etc/influxdb and do not forget to restart the server after every change in configuration file.**
```
- Telegraf Installation

```sh
- ***Download the latest .tar.gz file and extract it by using the following commands***

    cd $HOME
    wget https://dl.influxdata.com/telegraf/releases/telegraf-1.14.0_linux_amd64.tar.gz
    tar xf telegraf-1.14.0_linux_amd64.tar.gz

    **Start telegraph**

- go to the directory path and run the binary using below commands
    cd telegraf/usr/bin/
    ./telegraf --config ../../etc/telegraf/telegraf.conf

     **Note : If you want to give your own configuration then you can edit the telegraf.conf file which is in the path telegraf/etc/telegraf/telegraf.conf and do not forget to restart the server after every change in configuration file.**
```

```bash
git clone git@github.com:chris-remus/chainflow-vitwit.git
cd chainflow-vitwit
cp example.config.toml config.toml
```

- For **Get Missed Block Alerting** populate *missed_blocks_threshold* in `config.toml` with your desired threshold
- To **Get Validator Status Alerts** populate *alert_time1* and *alert_time2* in `config.toml` in the form of UTC(ex:"02:25PM")
- To **Get Block Difference Alerts** populate *block_diff_threshold* with your desired threshold
- For **Telegram Alerting** populate *chat_id* and *bot_token* in `config.toml` with your values
- For **Email Alerting** populate *token* and *to_email* (your email) in `config.toml` with your values

After populating config.toml -

```bash
- Build and run the bin file

    go build && ./chainflow-vitwit
```

```bash
docker build -t cfv .
docker run -d --name chainflow-vitwit cfv
```