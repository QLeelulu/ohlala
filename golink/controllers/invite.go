package controllers

import (
    "strings"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    //"strconv"
    //"fmt"
    //"github.com/QLeelulu/ohlala/golink"
)

type InviteResult struct {
	Result        bool
	Msg           string
	InviteUrl     string
}

type InviteViewModel struct {
	Title string
	RegisterInviteRemainCount int
}

/**
 * vote controller
 */
var _ = goku.Controller("invite").
    /**
     * 给指定的email发送邀请码
     */
    Get("email", func(ctx *goku.HttpContext) goku.ActionResulter {
	inviteModel := &InviteViewModel{"邀请", 0}
    var userId int64 = (ctx.Data["user"].(*models.User)).Id
	inviteModel.RegisterInviteRemainCount = models.RegisterInviteRemainCount(userId)
	return ctx.Render("/invite/show", inviteModel)

}).Filters(filters.NewRequireLoginFilter()).
    /**
     * 给指定的email发送邀请码
     */
    Post("email", func(ctx *goku.HttpContext) goku.ActionResulter {
    
    var userId int64 = (ctx.Data["user"].(*models.User)).Id
	if userId <= int64(0) {
		return ctx.Json(&InviteResult{false, "未登录", ""})
	}

	var strEmails string = ctx.Get("emails")

//fmt.Println(strEmails, "  ", userId)

	iCount := models.RegisterInviteRemainCount(userId)
	if strEmails == "" { //email为空代表获取邀请链接
		if iCount <= 0 {
			return ctx.Json(&InviteResult{false, "超出可以邀请的次数", ""})
		}
		inviteKey, err := models.CreateRegisterInviteWithoutEmail(userId)
		if err != nil {
			return ctx.Json(&InviteResult{false, "请求出错,请重试!", ""})
		}
		return ctx.Json(&InviteResult{true, "", "http://xxxx" + inviteKey})
	} else {
		arrEmails := strings.Split(strEmails, ";")
		if iCount < len(arrEmails) {
			return ctx.Json(&InviteResult{false, "超出可以邀请的次数", ""})
		}

		re, errReg := utils.GetEmailRegexp()
		if errReg != nil {
			return ctx.Json(&InviteResult{false, "请求出错,请重试!", ""})
		}
		for _, email := range arrEmails {
            if re.MatchString(email) == false {
                return ctx.Json(&InviteResult{false, "email格式不正确", ""})
            }
        }
		success, err := models.CreateRegisterInvite(userId, strEmails)
		if success == false {
			return ctx.Json(&InviteResult{false, "请求出错,请重试!", ""})
		}
		return ctx.Json(&InviteResult{true, "", ""})
	}

    return ctx.Json(&InviteResult{false, "请求出错,请重试!", ""})

}).Filters(filters.NewRequireLoginFilter())



