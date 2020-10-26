FROM golang:1.14.10

ARG WORKSPACE=/go/src/gola
COPY ./config/project/prod ${WORKSPACE}/config/project/prod
COPY ./config/schedule/prod ${WORKSPACE}/config/schedule/prod
COPY ./public ${WORKSPACE}/public
COPY ./docs ${WORKSPACE}/docs
WORKDIR ${WORKSPACE}

# 宣告環境變數
ENV PROJECT_ROOT=${WORKSPACE}
ENV TZ=Asia/Taipei

ARG version
# 最後放置執行檔，為了加快編譯速度
COPY ./gola${version} ${WORKSPACE}/gola

# 啟動服務
ENTRYPOINT [ "./gola" ]
