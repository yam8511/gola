###### Build Stage ######

FROM golang:1.14.10 AS builder

ARG WORKSPACE=/go/src/gola
WORKDIR ${WORKSPACE}
COPY . ${WORKSPACE}
RUN go build -mod vendor -o program


###### Final Stage ######

FROM golang:1.14.10

ARG WORKSPACE=/go/src/gola
WORKDIR ${WORKSPACE}

# 宣告環境變數
ENV PROJECT_ROOT=${WORKSPACE}
ENV TZ=Asia/Taipei

# 最後放置執行檔，為了加快編譯速度
COPY --from=builder ${WORKSPACE}/config ${WORKSPACE}/config
COPY --from=builder ${WORKSPACE}/docs ${WORKSPACE}/docs
COPY --from=builder ${WORKSPACE}/public ${WORKSPACE}/public
COPY --from=builder ${WORKSPACE}/storage ${WORKSPACE}/storage
COPY --from=builder ${WORKSPACE}/program ${WORKSPACE}/program

# 啟動服務
ENTRYPOINT [ "./program" ]
