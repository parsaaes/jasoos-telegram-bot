FROM golang:alpine As builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /jasoos-telegram-bot

FROM alpine:latest

LABEL maintainer="Parsa Eskandarnejad <parsa.eskandarnejad@gmail.com>"

WORKDIR /root/

COPY --from=builder /jasoos-telegram-bot .

ENTRYPOINT ["./jasoos-telegram-bot"]
