package filters

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/config"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
)

type ThirdPartyBindFilter struct{}

func (f *ThirdPartyBindFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {

    sessionIdBase, err := ctx.Request.Cookie(config.ThirdPartyCookieKey)
    if err != nil || len(sessionIdBase.Value) == 0 {
        ar = ctx.NotFound("no user binding context found.")
        return
    }
    ctx.Data["thirdPartySessionIdBase"] = sessionIdBase.Value

    profileSessionId := models.ThirdParty_GetThirdPartyProfileSessionId(sessionIdBase.Value)
    profile := models.ThirdParty_GetThirdPartyProfileFromSession(profileSessionId)

    if profile == nil {
        ar = ctx.NotFound("no user binding context found.")
        return
    }

    ctx.ViewData["profile"] = profile
    if len(profile.Email) > 0 {
        sensitiveInfoRemovedEmail := utils.GetSensitiveInfoRemovedEmail(profile.Email)
        ctx.ViewData["directCreateEmail"] = sensitiveInfoRemovedEmail
    }

    return
}

func (f *ThirdPartyBindFilter) OnActionExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (f *ThirdPartyBindFilter) OnResultExecuting(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (f *ThirdPartyBindFilter) OnResultExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func NewThirdPartyBindFilter() *ThirdPartyBindFilter {
    return &ThirdPartyBindFilter{}
}
