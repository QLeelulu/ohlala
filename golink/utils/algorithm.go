package utils

import (
    "github.com/QLeelulu/ohlala/golink"
    "math"
    "time"
)

//link 排序算法
func LinkSortAlgorithm(createTime time.Time, upVote int64, downVote int64) float64 {
	if upVote + downVote == 0 {
		return 0
	}
    var x int64 = upVote - downVote
    var y = 0.0
    var z int64 = 1
    switch {
    case x > 0:
        y = 1.0
        z = x
    case x < 0:
        y = -1.0
        z = -x
    }
    var ts = createTime.Sub(golink.SITERUNTIME_TIME)

    return math.Log10(float64(z)) + y*ts.Seconds()/golink.SCORETIMESTEMP
}

//link 排序算法
func CommentSortAlgorithm(upVote int64, downVote int64) float64 {

    n := float64(upVote + downVote)
    if n == 0.0 {
        return 0
	}
    z := 1.0 //1.0 = 85%, 1.6 = 95%
    phat := float64(upVote) / n

    return ( phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n) ) / (1+z*z/n)
}
