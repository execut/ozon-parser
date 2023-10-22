FROM golang:1.21.3

WORKDIR /usr/src/app

RUN apt-get update && apt-get -y upgrade && apt-get -y install gcc g++ ca-certificates xvfb chromium

COPY src/go.mod src/go.sum ./
RUN go mod download && go mod verify

COPY src .
RUN go build -v -o /usr/local/bin/ ./...

CMD ["ozon_parser", "-t=/data/token.txt", "positions", "-i=/data/words-for-positions.csv", "-o=/data/positions.csv"]
