# Build Stage
FROM golang:1.14-alpine AS builder

# 安裝基本工具
# RUN apk update && apk upgrade
# RUN apk add --no-cache bash git openssh gcc g++

# 複製原始碼
COPY . /app
WORKDIR /app

# 進行編譯
RUN go build -mod vendor -o gola


# Final Stage
FROM golang:1.14-alpine

COPY --from=builder /app/gola /app/gola
COPY --from=builder /app/config /app/config
COPY --from=builder /app/public /app/public
COPY --from=builder /app/storage /app/storage
WORKDIR /app

# 新增使用者
RUN adduser -D -u 1000 zuolar \
    # 調整 logs 的權限
    && chown -R zuolar:zuolar ./storage

# 宣告環境變數
ENV APP_ENV=
ENV APP_SITE=
ENV APP_ROOT=/app
ENV TZ=Asia/Taipei

ARG ENTER=server
ENV ENTER=${ENTER}

ARG BOT_TOKEN
ARG BOT_CHAT_ID
ENV BOT_TOKEN=${BOT_TOKEN}
ENV BOT_CHAT_ID=${BOT_CHAT_ID}

# 啟動服務
CMD [ "sh", "-c", "./gola $ENTER" ]

# 切換使用者
USER zuolar
