## Cosmos alerting bot

 -   **Cosmos alerting bot** will send alerts to your telegram account about your **validator status**(jailed or voting) and about **missed blocks**.

## Install Prerequisites
- **Go 13.x+**

### Getting Started

```bash
git clone https://github.com/Chainflow/cosmos-validator-mission-control.git
cd cosmos-validator-mission-control/alert-bot
cp example.config.toml config.toml
```
### Configure the following variables in `config.toml`

- *tg_chat_id*

    Telegram chat ID to receive Telegram alerts, required for Telegram alerting.
    
- *tg_bot_token*

    Telegram bot token, required for Telegram alerting. The bot should be added to the chat and should have send message permission.

- *alert_time1* and *alert_time2*

    These are for regular status updates. To receive validator status daily (twice), configure these parameters in the form of "02:25PM". The time here refers to UTC time.

- *val_operator_addr*

    Operator address of your validator which will be used to get staking, delegation and distribution rewards.

- *validator_hex_addr*

    Validator hex address useful to know about last proposed block, missed blocks and voting power.

- *lcd_endpoint*

    Address of your lcd client (ex: http://localhost:1317).

     **Note** : To start the lcd server for stargate node please do the following:
    
    - Go to `config/app.toml`
    - You can find below fields in that file
        ```bash 
           # Enable defines if the API server should be enabled.
           enable = false
        ```
    - Then make it `true` and restart the server.

- *external_rpc*

    External open RPC endpoint(secondary RPC other than your own validator). Useful to gather information like validator caught up, syncing and missed blocks etc.

After populating config.toml 

- Build and run the alerting bot using binary

```bash
$ go build -o cosmos-alert-bot && ./cosmos-alert-bot
```

- Run using docker

```bash
$ docker build -t alertbot .
$ docker run -d --name chainflow-vitwit alertbot
```