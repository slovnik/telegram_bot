FROM golang:latest
ARG bot_id
ARG api_url
ARG webhook_url
ADD . /go/src/github.com/slovnik/telegram_bot
WORKDIR /go/src/github.com/slovnik/telegram_bot
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go install github.com/slovnik/telegram_bot
ENV SLOVNIK_BOT_ID $bot_id
ENV SLOVNIK_API_URL $api_url
ENV SLOVNIK_WEBHOOK_URL $webhook_url
ENTRYPOINT /go/bin/telegram_bot
EXPOSE 8080
