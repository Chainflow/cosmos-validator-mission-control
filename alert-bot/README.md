## Validator monitoring and alerting tool

### Getting Started

```bash
git clone https://github.com/Chainflow/cosmos-validator-mission-control.git
git fetch && git checkout akash_alert_bot
cd cosmos-validator-mission-control/alert-bot
cp example.config.toml config.toml
```
- *val_operator_addr*
    Operator address of your validator which will be used to get staking, delegation and distribution rewards.

- *account_addr* 

    Your validator account address which will be used to get account informtion etc.

- *validator_hex_addr*

    Validator hex address useful to know about last proposed block, missed blocks and voting power.

- *lcd_endpoint*

    Address of your lcd client (ex: http://localhost:1317).

- *external_rpc*

    External open RPC endpoint(secondary RPC other than your own validator). Useful to gather information like validator caught up, syncing and missed blocks etc.

- *alert_time1* and *alert_time2*

    These are for regular status updates. To receive validator status daily (twice), configure these parameters in the form of "02:25PM". The time here refers to UTC time.

- *tg_chat_id*

    Telegram chat ID to receive Telegram alerts, required for Telegram alerting.
    
- *tg_bot_token*

    Telegram bot token, required for Telegram alerting. The bot should be added to the chat and should have send message permission.

- *email_address*

    E-mail address to receive mail notifications, required for e-mail alerting.

- *sendgrid_token*

    Sendgrid mail service api token, required for e-mail alerting.

-  *validator_name*

    Give your validator name to get validator name in telegram and email notifications.


After populating config.toml -

### 1. Build and run

```bash
go build && ./alert-bot
```

### 2. Run using docker
```bash
docker build -t cfv .
docker run -d --name chainflow-vitwit cfv
```
### 3. Run using systemd services

If you wish to run systemd service for alert-bot make sure to follow these steps.
```bash
 i. Copy config.toml file in user directory.
ii. Do `go build` in alert-bot directory to generate binary
iii. In ExecStart=<give path of alert-bot binary> (ex: ExecStart=go/src/github.com/Chainflow/cosmos-validator-mission-control/alert-bot/alert-bot)
```

```bash
echo "[Unit]
Description=Alert-bot System Service
After=network-online.target

[Service]
User=<your_user>
ExecStart=cd <alert-bot-path> && go run main.go
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" > alertbot.service
```
Enable and start the systemd service using:

```bash
sudo mv alertbot.service /lib/systemd/system/alertbot.service
sudo -S systemctl daemon-reload
sudo -S systemctl start alertbot
```

### You can run akash rest server by using systemd service

Fetch akashctl location path and use it in next step.

```bash
which akashctl
```

Paste in the following

```bash
echo "[Unit]
Description=Akash System Service
After=network-online.target

[Service]
User=<your_user>
ExecStart=<akashctl-path> rest-server
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" > akashctl.service
```

Enable and start the systemd service using:

```bash
sudo mv akashctl1.service /lib/systemd/system/akashctl.service
sudo -S systemctl daemon-reload
sudo -S systemctl start akashctl
```