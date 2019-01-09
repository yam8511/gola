# Build Stage
FROM golang:1.11-alpine AS builder

# 安裝基本工具
# RUN apk update && apk upgrade
# RUN apk add --no-cache bash git openssh gcc g++

# 複製原始碼
COPY . /go/src/gola
WORKDIR /go/src/gola

# 進行編譯
RUN go build -o gola


# Final Stage
FROM alpine

ARG APP_ENV=docker
COPY --from=builder /usr/local/go/lib /usr/local/go/lib
COPY --from=builder /go/src/gola/gola /app/gola
COPY ./config /app/config
COPY ./public /app/public
COPY ./storage /app/storage
WORKDIR /app

# 新增使用者
RUN adduser -D -u 1000 zuolar \
    # 調整 logs 的權限
    && chown -R zuolar:zuolar ./storage

# 宣告環境變數
ENV GIN_MODE=debug APP_ENV=docker TZ=Asia/Taipei

# 啟動服務
ENTRYPOINT [ "./gola" ]

# 切換使用者
USER zuolar
