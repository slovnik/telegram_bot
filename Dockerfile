FROM golang:latest
ARG bot_id
ADD . /go/src/github.com/slovnik/telegram_bot
WORKDIR /go/src/github.com/slovnik/telegram_bot
RUN go get github.com/tools/godep
RUN godep restore
RUN go install github.com/slovnik/telegram_bot
ENV SLOVNIK_BOT_ID $bot_id
ENTRYPOINT /go/bin/telegram_bot
