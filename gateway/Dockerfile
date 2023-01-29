FROM golang:1.19.5-alpine3.16

WORKDIR /app

COPY go.mod .

RUN go mod download

ADD cmd/config.yaml /app

COPY cmd/ssl/ /app/ssl/

COPY . .

RUN go build -o gateway-svc cmd/main.go 

EXPOSE 8000

CMD [ "/app/gateway-svc" ]