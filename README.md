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

- *voting_power_threshold*

    Configure the threshold to receive alert when the voting power reaches or drops below of the threshold which has been given.

- *num_peers_threshold*

    Configure the threshold to get an alert if the no.of connected peers falls below of the threshold.

After populating config.toml, build and run the monitoring binary

```bash
$ go build -o chain-monit && ./chain-monit
```

```bash
$ docker build -t cfv .
$ docker run -d --name chain-monit cfv
```

In grafana there will be two types of dhasboard 
```bash
i. Validator Monitoring Metrics (These are the metrics which we have calculated and stored in influxdb)
ii. System Metrics (These are related to system configuration and all and which comes from telegraf)
```

**List of validator monitoring metrics**

- Validator Details :  Which displays details of a validator like moniker, website, keybase and details.
- Gaiad Status :  Displays the chain is running or not in the from of UP and DOWN.
- Validator Status :  Displays about the validator health like Voting if the validator is in active state or else Jailed.
- Gaiad Version : Displays the version of gaia.
- Validator Caught Up : Displays whether the validator node has been synced to the network or not.
- Block Time Difference : Displays the time difference between previous block and current block syncing.
- Current Block Height -  Validator : Displays the current block height of validator.
- Latest Block Height - Network : Displays the latest block height of a network.
- Height Difference : Displays the difference between heights of validator current block height and network latest block height.
- Last Missed Block Range : Displays the continuous missed blocks range based on the missed block thershold given in the config.toml
- Blocks Missed In last 48h : Displays the count of blocks missed by a validator in last 48 hours.
- Unconfirmed Txns : Displays the number of uncofirmed transactions in that node.
- No.of Peers : Displays the total number of peers connected in a network.
- Peer Address : Displays the addresses of connected peers.
- Latency : Displays the latency of connected peers in the form of graph.
- Validator Fee : Displays the commission rate of a validator as fee.
- Max Change Rate : Displays the max change rate of a validator commission.
- Max Rate : Displays the max rate of a validator commission.
- Voting Power : Displays the voting power of a validator which has given in config.
- Self Delegation Balance : Displays delegation balance of your validator.
- Current Balance : Displays the account balance of your validator.
- Unclaimed Rewards : Displays the current rewards amount the validator.
- Last proposed Block Heigt : Displays height of the last proposed block if it has been proposed by your validator.
- Last Proposed Block Time : Displays time of the last proposed block which has been proposed be your validator.
- Voting Period Proposals : Displays list of the proposals which are in voting period.
- Deposit Period Proposals : Displays list of the proposals which are in deposit period.
- Completed Proposals : Displays list of the proposals which are completed it might be passed or rejected.

```bash
Note: Above mentioned metrics will be calculated and displayed according to the validator address you will be populating in config.toml
```

**System Monitoring Metrics**
-  For this you can refer `telgraf.conf` file for system monitoring metrics.You can just replace it with your original telegraf.conf file which have been located at /telegraf/etc/telegraf
 

 **Alerting (Telegram and Email)**

 - Alert about validator node sync status.
 - Alert about missed blocks when the missed blocks count reaches or exceeded **missed_blocks_threshold** which has been given by user in *config.toml*
 - Alert about no.of peers when the count falls below of **num_peers_threshold** which has been given by user in *config.toml*
- Alert about the block difference between network and validator reaches or exceeds the **block_diff_threshold** which has been given by user in *config.toml*
- Alert about the gaiad status whether it's running on your validator instance or not.
- Alert about a new proposal
- Alert about the proposal if it's moved to voting period, passed or rejected.
- Alert about voting period proposals if the voting end time is less than or equal to 24 hours and also if the validator didn't vote on propoal yet.
- Alert about validator health whether it's voting or jailed. You can get alerts twice a day based on the time you have configured **alert_time1** and **alert_time2** in *config.toml*
- Alert about the voting power of your validator when it reaches or drops below of the **voting_power_threshold** which has been given by user in *config.toml*
