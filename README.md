# gola

---

## 啟動本機開發的資料庫

    ```shell
    # 啟動開發用的DB
    $ docker-compose up -d db
    ```

    ps. 啟動DB之後，稍等一下DB程序啟動完畢
    ps. 看http://127.0.0.1:8080，手動調整DB的編碼，改成 utf8_general_ci

---

## 編譯程式

- 方法1. 啟動/關閉開發用的容器

    ```shell
    # 容器編譯程式
    $ docker-compose up build
    # 容器啟動 server
    $ docker-compose up web-dev
    # 關閉所有容器
    $ docker-compose down
    ```

- 方法2. 手動編譯檔案

    ```shell
    # 使用 golang v1.11 版本
    # 啟用 go mod (內建)
    $ export GO111MODULE=on
    # 編譯程式
    $ go build -o gola
    # 啟動 server
    $ APP_ENV=local ./gola server
    ```

- 方法3. 使用容器編譯，手動啟動server

    ```shell
    # 編譯程式
    $ docker-compose up build
    # 啟動 server
    $ APP_ENV=local ./gola server
    ```

---

## [查看頁面](http://127.0.0.1:8000)

---

## 部署程式

```shell
APP_ENV=prod docker-compose build build-image
APP_ENV=prod docker-compose up -d web cronjob
```
