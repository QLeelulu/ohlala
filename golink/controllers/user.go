package controllers

import (
    "crypto/md5"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "html/template"
    "net/http"
    "strconv"
    "strings"
    "time"
)

/**
 * form
 */
func createLoginForm() *form.Form {
    // defined the field
    email := form.NewEmailField("email", "Email", true).
        Error("invalid", "Email地址错误").
        Error("require", "Email地址必须填写").Field()

    pwd := form.NewCharField("pwd", "密码", true).Min(6).Max(30).
        Error("required", "密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    remeber_me := form.NewCharField("remeber_me", "记住我", false)

    // add the fields to a form
    form := form.NewForm(email, pwd, remeber_me)
    return form
}

func createRegForm() *form.Form {

    key := form.NewCharField("key", "注册码", true).
        Error("required", "注册码必须填写").Field()

    email := form.NewEmailField("email", "Email", true).
        Error("invalid", "Email地址错误").
        Error("required", "Email地址必须填写").Field()

    pwd := form.NewCharField("pwd", "密码", true).Min(6).Max(30).
        Error("required", "密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    name := form.NewCharField("name", "昵称", true).Min(2).Max(15).
        Error("required", "昵称必须填写").
        Error("range", "昵称长度必须在{0}到{1}之间").Field()

    repwd := form.NewCharField("repwd", "确认密码", true).Min(6).Max(30).
        Error("required", "确认密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(key, email, name, pwd, repwd) 
    return form
}

// 更新基本信息
func createBaseInfoForm() *form.Form {
    description := form.NewTextField("description", "自我介绍", false).Max(100).
        Error("max-length", "自我介绍的字数不能多于{0}个").Field()

    name := form.NewCharField("name", "用户名", true).Min(2).Max(15).
        Error("required", "用户名必须填写").
        Error("range", "用户名长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(description, name)
    return form
}

// 更新密码
func createUpdatePwdForm() *form.Form {
    oldPwd := form.NewCharField("old-pwd", "原密码", true).Min(6).Max(30).
        Error("required", "原密码必须填写").
        Error("range", "原密码长度必须在{0}到{1}之间").Field()

    newPwd := form.NewCharField("new-pwd", "新密码", true).Min(6).Max(30).
        Error("required", "新密码必须填写").
        Error("range", "新密码长度必须在{0}到{1}之间").Field()

    newPwd2 := form.NewCharField("new-pwd2", "确认密码", true).Min(6).Max(30).
        Error("required", "确认密码必须填写").
        Error("range", "确认密码长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(oldPwd, newPwd, newPwd2)
    return form
}

/**
 * Controller: user
 */
var _ = goku.Controller("user").

    /**
     * login view
     */
    Get("login", func(ctx *goku.HttpContext) goku.ActionResulter {

    if u, ok := ctx.Data["user"]; ok && u != nil {
        return ctx.Redirect("/")
    }
    ctx.ViewData["query"] = template.URL(ctx.Request.URL.RawQuery)
    return ctx.View(nil)
}).

    /**
     * reg view
     */
    Get("reg", func(ctx *goku.HttpContext) goku.ActionResulter {

    if u, ok := ctx.Data["user"]; ok && u != nil {
        return ctx.Redirect("/")
    }
    ctx.ViewData["query"] = template.URL(ctx.Request.URL.RawQuery)
    ctx.ViewData["code"] = ctx.Get("key")
    return ctx.Render("login", nil)
}).

    /**
     * logout
     */
    Get("logout", func(ctx *goku.HttpContext) goku.ActionResulter {

    redisClient := models.GetRedis()
    defer redisClient.Quit()
    redisClient.Del("_glut")
    c := &http.Cookie{
        Name:    "_glut",
        Expires: time.Now().Add(-10 * time.Second),
        Path:    "/",
    }
    ctx.SetCookie(c)
    return ctx.Redirect("/")
}).

    /**
     * login
     */
    Post("login", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := createLoginForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        email := strings.ToLower(m["email"].(string))
        // 检查密码是否正确
        userId := models.User_CheckPwd(email, m["pwd"].(string))
        if userId > 0 {
            now := time.Now()
            h := md5.New()
            h.Write([]byte(fmt.Sprintf("%v-%v", email, now.Unix())))
            ticket := fmt.Sprintf("%x_%v", h.Sum(nil), now.Unix())
            var expires time.Time
            if m["remeber_me"] == "1" {
                expires = now.Add(365 * 24 * time.Hour)
            } else {
                expires = now.Add(48 * time.Hour)
            }
            redisClient := models.GetRedis()
            defer redisClient.Quit()
            err := redisClient.Set(ticket, userId)
            if err != nil {
                goku.Logger().Errorln(err.Error())
                errorMsgs = append(errorMsgs, "登陆票据服务器出错，请重试")
            } else {
                _, err = redisClient.Expireat(ticket, expires.Unix())
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                    errorMsgs = append(errorMsgs, "登陆票据服务器出错，请重试")
                }
                c := &http.Cookie{
                    Name:    "_glut",
                    Value:   ticket,
                    Expires: expires,
                    //Domain:     ".godev.local", // edit (or omit)
                    Path:     "/", // ^ ditto
                    HttpOnly: true,
                }
                ctx.SetCookie(c)
            }
        } else {
            errorMsgs = append(errorMsgs, "账号密码不正确，请改正")
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        // ctx.ViewData["loginSuccess"] = true
        returnUrl := ctx.Request.FormValue("returnurl")
        if returnUrl == "" {
            returnUrl = "/"
        }
        return ctx.Redirect(returnUrl)
    } else {
        ctx.ViewData["loginErrors"] = errorMsgs
        ctx.ViewData["loginValues"] = f.Values()
    }

    return ctx.Render("login", nil)
}).

    /**
     * submit reg
     */
    Post("reg", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := createRegForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        if m["pwd"] == m["repwd"] {

            // 检查email地址是否已经注册
            emailExist := models.User_IsEmailExist(m["email"].(string))
            userExist := models.User_IsUserExist(m["name"].(string))
			regKey := models.VerifyInviteKey(m["key"].(string))
            if !emailExist && !userExist && regKey != nil { 
                m["pwd"] = utils.PasswordHash(m["pwd"].(string))
                delete(m, "repwd")
				delete(m, "key")
                m["create_time"] = time.Now()
                _, err := models.User_SaveMap(m)
				if err == nil {
					models.UpdateIsRegister(regKey)
				} else {
                    errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
                    goku.Logger().Errorln(err)
                }
            } else {
                if userExist {
                    errorMsgs = append(errorMsgs, "用户名已经被注册，请换一个")
                }
                if emailExist {
                    errorMsgs = append(errorMsgs, "Email地址已经被注册，请换一个")
                }
				if regKey == nil {
                    errorMsgs = append(errorMsgs, "注册码不正确，可能已经被注册或过期")
				}
            }
        } else {
            errorMsgs = append(errorMsgs, "两次输入的密码不一样")
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        ctx.ViewData["regSuccess"] = true
    } else {
        ctx.ViewData["regErrors"] = errorMsgs
        ctx.ViewData["regValues"] = f.Values()
    }

    return ctx.Render("login", nil)
}).

    /**
     * follow somebody
     */
    Post("follow", func(ctx *goku.HttpContext) goku.ActionResulter {

    followId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.User_Follow(ctx.Data["user"].(*models.User).Id, followId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * follow somebody
     */
    Post("unfollow", func(ctx *goku.HttpContext) goku.ActionResulter {

    followId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.User_UnFollow(ctx.Data["user"].(*models.User).Id, followId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 查看用户信息页
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {

    userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    user := models.User_GetById(userId)

    if user == nil {
        ctx.ViewData["errorMsg"] = "用户不存在"
        return ctx.Render("error", nil)
    }

    links := models.Link_ByUser(user.Id, 1, 10)
    friends, _ := models.UserFollow_Friends(user.Id, 1, 12)
    followers, _ := models.UserFollow_Followers(user.Id, 1, 12)
    followTopics, _ := models.User_GetFollowTopics(user.Id, 1, 12)

    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["Friends"] = friends
    ctx.ViewData["Followers"] = followers
    ctx.ViewData["FollowTopics"] = followTopics
    return ctx.View(models.User_ToVUser(user, ctx))

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 查看用户设置页
     */
    Get("setting", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    return ctx.View(models.User_ToVUser(user, ctx))

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 更新用户基本信息
     */
    Post("update-base", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    f := createBaseInfoForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        // 检查用户名是否已经注册
        userExist := false
        if strings.ToLower(m["name"].(string)) != strings.ToLower(user.Name) {
            userExist = models.User_IsUserExist(m["name"].(string))
        }
        if userExist {
            errorMsgs = append(errorMsgs, "用户名已经被注册，请换一个")
        } else {
            _, err := models.User_Update(user.Id, m)
            if err != nil {
                errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
                goku.Logger().Errorln(err)
            } else {
                user = models.User_GetById(user.Id)
            }
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        ctx.ViewData["updateBaseSuccess"] = true
    } else {
        ctx.ViewData["updateBaseErrors"] = errorMsgs
        v := f.Values()
        user.Name = v["name"]
        user.Description = v["description"]
    }

    return ctx.Render("setting", user)

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 更新用户密码
     */
    Post("change-pwd", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    f := createUpdatePwdForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        if m["new-pwd"] == m["new-pwd"] {
            // 检查原密码是否正确
            if utils.PasswordHash(m["old-pwd"].(string)) == user.Pwd {
                saveMap := map[string]interface{}{"pwd": utils.PasswordHash(m["new-pwd"].(string))}
                _, err := models.User_Update(user.Id, saveMap)
                if err != nil {
                    errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
                    goku.Logger().Errorln(err)
                }
            } else {
                errorMsgs = append(errorMsgs, "原密码不正确，请重新输入")
            }
        } else {
            errorMsgs = append(errorMsgs, "两次输入的新密码不一致")
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        ctx.ViewData["updatePwdSuccess"] = true
    } else {
        ctx.ViewData["updatePwdErrors"] = errorMsgs
    }

    return ctx.Render("setting", user)

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 获取用户信息
     * 用于浮动层
     */
    Get("pbox-info", func(ctx *goku.HttpContext) goku.ActionResulter {

    userId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    user := models.User_GetById(userId)

    if user != nil {
        return ctx.RenderPartial("pop-info", models.User_ToVUser(user, ctx))
    }
    return ctx.Html("")

}).Filters(filters.NewAjaxFilter())
