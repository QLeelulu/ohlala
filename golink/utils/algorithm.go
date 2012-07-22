package algorithm

import (
    "math"
    "time"
)

var startTime = time.Date(2012, time.July, 21, 23, 20, 25, 0, time.UTC)
//reddit 排序算法
func RedditSortAlgorithm(createTime time.Time, upVote int64, downVote int64) float64 {
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
    var ts = createTime.Sub(startTime)

    return math.Log10(float64(z)) + y * ts.Seconds()/45000
}

