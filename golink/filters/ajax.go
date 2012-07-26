package filters

import (
    "github.com/QLeelulu/goku"
)

// 检查是否为AJAX请求
type AjaxFilter struct {
}

func (tf *AjaxFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
    if !ctx.IsAjax() {
        return ctx.Raw("Not AJAX"), nil
    }
    return nil, nil
}
func (tf *AjaxFilter) OnActionExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (tf *AjaxFilter) OnResultExecuting(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (tf *AjaxFilter) OnResultExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

var ajaxFilter = new(AjaxFilter)

func NewAjaxFilter() *AjaxFilter {
    return ajaxFilter
}
