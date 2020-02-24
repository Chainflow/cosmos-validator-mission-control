FROM golang:1.13.1

COPY . /app/chainflow-vitwit

WORKDIR /app/chainflow-vitwit

RUN go mod download

CMD ["go", "run", "main.go"]