# GoLa

- Werewolf 狼人殺
- Crimnal Dance 犯人在跳舞

---

## 啟動資料庫

```shell
make setup
```

---

## 執行程式

```shell
# 編譯程式
go build -o gola .

# 啟動伺服器
./gola server

# 啟動排程
./gola schedule

# 執行指令
# ./gola run [自定義指令]
./gola run demo
```

---

## 查看頁面

- [GraphQL Demo](http://127.0.0.1:8000/graphql)
- [狼人殺](http://127.0.0.1:8000/wf)
- [犯人在跳舞](http://127.0.0.1:8000/cd)

---

## 部署程式

```shell
docker-compose build web
docker-compose push web
# 需登入 heroku
heroku container:release web -a golar
```
