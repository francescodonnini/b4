FROM golang:1.21-alpine

RUN apk update

WORKDIR /app

RUN apk add iproute2 iputils

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /b4

EXPOSE 5050

ENTRYPOINT ["sh", "setup.sh"]
