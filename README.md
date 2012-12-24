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

### 项目结构说明 ###

```bash
.
├── app.go web服务启动文件
├── golink
│   ├── config.go  配置文件
│   ├── controllers  
│   │   ├── admin 管理后台
│   │   │   ├── base.go
│   │   │   ├── comment.go
│   │   │   ├── index.go
│   │   │   ├── link.go
│   │   │   ├── topic.go
│   │   │   └── user.go
│   │   ├── api.go  
│   │   ├── base.go 一些controller相关的公用函数
│   │   ├── comment.go
│   │   ├── discover.go
│   │   ├── home.go
│   │   ├── host.go
│   │   ├── invite.go
│   │   ├── link.go
│   │   ├── topic.go
│   │   ├── user.go
│   │   ├── user_reg_login.go  用户登陆、注册
│   │   ├── user_setting.go    用户设置
│   │   └── vote.go
│   ├── db
│   │   └── link.sql   数据脚本
│   ├── filters        Controller/Action Filter
│   │   ├── ajax.go
│   │   └── require_login.go
│   ├── forms          Form表单验证
│   │   └── forms.go
│   ├── global_viewdata.go
│   ├── middlewares                中间件，对所有请求做统一处理
│   │   ├── confess.go             前期推广不容易啊！
│   │   └── util-middleware.go     例如判断用户是否登陆等一些常用中间件
│   ├── models
│   │   ├── admin_comment.go
│   │   ├── admin_link.go
│   │   ├── base.go
│   │   ├── comment.go
│   │   ├── comment_for_user.go      用户收到的评论
│   │   ├── comment_sort.go          查看Link时评论列表排序
│   │   ├── invite.go
│   │   ├── link.go
│   │   ├── link_for_home.go         
│   │   ├── link_for_host.go         
│   │   ├── link_for_topic.go
│   │   ├── link_for_user.go
│   │   ├── link_support_record.go
│   │   ├── remind.go                新评论、关注提醒
│   │   ├── remind_test.go
│   │   ├── topic.go
│   │   ├── user.go
│   │   ├── user_follow.go
│   │   └── vote.go
│   ├── route.go       Url路由配置
│   ├── static
│   │   ├── css
│   │   ├── ico
│   │   ├── img
│   │   └── js
│   │       ├── comment.js
│   │       ├── invite.js
│   │       ├── main.js
│   │       ├── seajs-lib       Seajs模块
│   │       ├── topic.js
│   │       ├── user-page.js
│   │       └── util.js
│   ├── utils                   公用帮助类
│   │   ├── algorithm.go
│   │   ├── genetic_key.go
│   │   ├── sina_oauth.go
│   │   ├── sina_weibo.go
│   │   ├── utils.go
│   │   └── utils_test.go
│   └── views                   视图(模板)
│       ├── _golink_admin       管理后台
│       └── shared
│           └── layout.html     主布局模板
├── golink.conf.sample          JSON配置文件
├── push-to-topic-and-home.go   推送到主题和首页的后台任务
├── push-to-user.go             把用户关注的内容推送给用户，后台任务
└── send-invite-email.go        发送邀请注册Email，后台任务
```
