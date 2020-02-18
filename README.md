## Validator monitoring and alerting tool

### Getting Started

```bash
git clone git@github.com:chris-remus/chainflow-vitwit.git
cd chainflow-vitwit
cp example.config.toml config.toml
```

- For **Telegram Alerting** populate *chat_id* and *bot_token* in `config.toml` with your values
- For **Email Alerting** populate *token* and *to_email* (your email) in `config.toml` with your values

After populating config.toml -
```bash
docker build -t cfv .
docker run -d --name chainflow-vitwit cfv
```