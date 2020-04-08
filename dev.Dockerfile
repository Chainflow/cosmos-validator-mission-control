FROM golang:1.13.1

COPY . /app/chain-monit

WORKDIR /app/chain-monit

RUN go mod download

CMD ["go", "run", "main.go"]