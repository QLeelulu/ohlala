package golink

// 排序
const (
    ORDER_TYPE_HOT  = "hot"  // 热门，按时间与票数算得的分数排序
    ORDER_TYPE_HOTC = "hotc" // 热评，按评论数排序
    ORDER_TYPE_TIME = "time" // 最新，按时间排序
    ORDER_TYPE_CTVL = "ctvl" // 争议，按 顶/踩 数排序
    ORDER_TYPE_VOTE = "vote" // 得分，按投票数排序
)

var ORDER_TYPE_MAP map[string]string = map[string]string{
    ORDER_TYPE_HOT:  "hot",
    ORDER_TYPE_HOTC: "hotc",
    ORDER_TYPE_TIME: "time",
    ORDER_TYPE_CTVL: "ctvl",
    ORDER_TYPE_VOTE: "vote",
}
