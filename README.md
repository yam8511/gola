# GoLa

- Werewolf 狼人殺
- Crimnal Dance 犯人在跳舞

---

## 編譯程式

```shell
# 編譯程式
$ go run main.go server
```

---

## 查看頁面

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
