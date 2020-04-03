## Validator monitoring and alerting tool

### Getting Started

```bash
git clone git@github.com:chris-remus/chainflow-vitwit.git
cd chainflow-vitwit
cp example.config.toml config.toml
```

- For **Get Missed Block Alerting** populate *missed_blocks_threshold* in `config.toml` with your desired threshold
- To **Get Validator Status Alerts** populate *alert_time1* and *alert_time2* in `config.toml` in the form of UTC(ex:"02:25PM")
- For **Telegram Alerting** populate *chat_id* and *bot_token* in `config.toml` with your values
- For **Email Alerting** populate *token* and *to_email* (your email) in `config.toml` with your values

After populating config.toml -
```bash
docker build -t cfv .
docker run -d --name chainflow-vitwit cfv
```