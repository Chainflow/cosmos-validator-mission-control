FROM golang:1.13.1

COPY . /app/chainflow-vitwit/alert-bot

WORKDIR /app/chainflow-vitwit/alert-bot

RUN go mod download

CMD ["go", "run", "main.go"]