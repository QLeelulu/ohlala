## 这是哪里 ##

这里是 [觅链](http://milnk.com) 的源码。

### 觅链[milnk] 是什么 ###

[觅链](http://milnk.com)是一个具有社交媒体属性的链接分享与评论平台。

### 用了哪些技术 ###

###### 后端 ######

-   [golang](http://golang.org/)
-   [goku](https://github.com/QLeelulu/goku)
-   mysql
-   [redis](http://redis.io/)

###### 前端 ######

-   jquery
-   [seajs](http://seajs.org/)
-   [bootstrap](http://twitter.github.com/bootstrap/)

### 怎样运行 ###

先建数据库： [db/link.sql](https://github.com/QLeelulu/ohlala/blob/master/golink/db/link.sql)

修改 `golink/config.go` 的相关配置，然后执行：

```bash
$go run app.go
```

