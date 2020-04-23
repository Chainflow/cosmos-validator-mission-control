# Validator monitoring and alerting tool

## Prerequisites
- **Go 13.x+**

**Install Grafana on ubuntu**

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

Grafana will be running on port :3000 (ex:: https://localhost:3000)
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

The default port that runs the InfluxDB HTTP service is :8086
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

By default telegraf does not expose any ports.
```

## Get the code
```bash
$ git clone git@github.com:chris-remus/chainflow-vitwit.git
$ cd chainflow-vitwit
$ git fetch && git checkout develop
$ cp example.config.toml config.toml
```

`config.toml` has following configurations
- *tg_chat_id*

    Telegram chat ID to receive telegram alerts
- *tg_bot_token*

    Telegram bot token. The bot should be added to the chat and it should have send message permission

- *email_address*

    E-mail address to receive mail notifications

- *sendgrid_token*

    Sendgrid mail service api token.
- *missed_blocks_threshold*

    Configure the threshold to receive  **Missed Block Alerting**
- *block_diff_threshold*

    An integer value to receive **Block difference alerts**

- *alert_time1* and *alert_time2*

    These are for regular status updates. To receive validator status daily (twice), configure these parameters in the form of "02:25PM". The time here refers to UTC time.

- *voting_power_threshold*

    Configure the threshold to receive alert when the voting power reaches or drops below of the threshold given.

- *num_peers_threshold*

    Configure the threshold to get an alert if the no.of connected peers falls below the threshold.

- *enable_telegram_alerts*

    Configure **yes** if you wish to get telegram alerts otherwise make it **no** .

- *enable_email_alerts*

    Configure **yes** if you wish to get email alerts otherwise make it **no** .

After populating config.toml, build and run the monitoring binary

```bash
$ go build -o chain-monit && ./chain-monit
```

Run using docker
```bash
$ docker build -t cfv .
$ docker run -d --name chain-monit cfv
```

In grafana there will be two types of dashboard 
```bash
i. Validator Monitoring Metrics (These are the metrics which we have calculated and stored in influxdb)
ii. System Metrics (These are related to system configuration which comes from telegraf)
```

**List of validator monitoring metrics**

- Validator Details :  Displays the details of a validator like moniker, website, keybase identity and details.
- Gaiad Status :  Displays whether the gaiad is running or not in the from of UP and DOWN.
- Validator Status :  Displays the validator health. Shows Voting if the validator is in active state or else Jailed.
- Gaiad Version : Displays the version of gaia currently running.
- Validator Caught Up : Displays whether the validator node is in sync with the network or not.
- Block Time Difference : Displays the time difference between previous block and current block.
- Current Block Height :  Validator : Displays the current block height committed by the validator.
- Latest Block Height : Network : Displays the latest block height of a network.
- Height Difference : Displays the difference between heights of validator current block height and network latest block height.
- Last Missed Block Range : Displays the continuous missed blocks range based on the threshold given in the config.toml
- Blocks Missed In last 48h : Displays the count of blocks missed by the validator in last 48 hours.
- Unconfirmed Txns : Displays the number of unconfirmed transactions on that node.
- No.of Peers : Displays the total number of peers connected to the validator.
- Peer Address : Displays the ip addresses of connected peers.
- Latency : Displays the latency of connected peers with respect to the validator.
- Validator Fee : Displays the commission rate of the validator.
- Max Change Rate : Displays the max change rate of the commission.
- Max Rate : Displays the max rate of the commission.
- Voting Power : Displays the voting power of the validator.
- Self Delegation Balance : Displays delegation balance of the validator.
- Current Balance : Displays the account balance of the validator.
- Unclaimed Rewards : Displays the current unclaimed rewards amount of the validator.
- Last proposed Block Height : Displays height of the last block proposed by the validator.
- Last Proposed Block Time : Displays the time of the last block proposed by the validator.
- Voting Period Proposals : Displays the list of the proposals which are currently in voting period.
- Deposit Period Proposals : Displays the list of the proposals which are currently in deposit period.
- Completed Proposals : Displays the list of the proposals which are completed with their status as passed or rejected.

```bash
Note: Above mentioned metrics will be calculated and displayed according to the validator address which will be configured in config.toml
```

**System Monitoring Metrics**

- Telegraf is a daemon that can run on any server and collect a wide variety of metrics from the system (cpu, memory, swap, etc.), common services (mysql, redis, postgres, etc.). It was originally built as a metric-gathering agent for InfluxDB, but has recently evolved to output metrics to other data sinks as well, such as Kafka, Datadog, and OpenTSDB.

-  For system monitoring metrics you can refer `telgraf.conf` file. You can just replace it with your original telegraf.conf file which will be located at /telegraf/etc/telegraf
 

 **Alerting** - Telegram and email are the platforms used to send alerts to the end user. Telegram was chosen over other platforms like Discord etc. as the majority of validators prefer to use Telegram for their personal alerting mechanisms as well. Telegarm alerts will be sent to chat_id present in config.toml and Email alerts will be sent to the email id present in config.toml.

 - Alert about validator node sync status.
 - Alert when missed blocks when the missed blocks count reaches or exceedes **missed_blocks_threshold** which is user configured in *config.toml*
 - Alert when no.of peers when the count falls below of **num_peers_threshold** which is user configured in *config.toml*
- Alert when the block difference between network and validator reaches or exceeds the **block_diff_threshold** which is user configured in *config.toml*
- Alert when the gaiad status is not running on validator instance.
- Alert when a new proposal is created.
- Alert when the proposal is moved to voting period, passed or rejected.
- Alert when voting period proposals is due in less than 24 hours and also if the validator didn't vote on proposal yet.
- Alert about validator health whether it's voting or jailed. You can get alerts twice a day based on the time you have configured **alert_time1** and **alert_time2** in *config.toml*
- Alert when the voting power of your validator drops below **voting_power_threshold** which is user configured in *config.toml*


**Instructions to setup the dashboards in grafana**

*Login*
- Open your web browser and go to http://localhost:3000/.  3000 is the default HTTP port that Grafana listens to if you haven’t configured a different port.
- If you are a first time user type admin for the username and password in the login page.
- You can change the password after login.

*Import the dashboards*
- To import the json file of the **validator monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the validator_monitoring_metrics.json present in the grafana_template folder. 

- Select the datasources and click on import.

- To import **system monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the system monitoring metrics.json present in the grafana_template folder

- While creating this dashboard if you face any issues at valueset, change it to empty and then click on import by selecting the datasources.

- *For more info about grafana dashboard imports you can refer https://grafana.com/docs/grafana/latest/reference/export_import/*
