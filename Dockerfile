FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && mkdir /log \
    && go build


FROM alpine:latest as runner
ARG env
RUN mkdir -p /app/config \
    && mkdir -p /app/logs

WORKDIR /app

COPY --from=builder /app/openai-chatGPT-server /app/
COPY --from=builder /app/config/dev_config.toml /app/config/

EXPOSE 8090

ENTRYPOINT /app/openai-chatGPT-server -env dev >> /app/logs/log.log 2>&1
