package utils

import (
    "github.com/sdegutis/go.assert"
    "testing"
)

func TestIsSpider(t *testing.T) {
    ua := "Chrome ohlala xsd sf sdf "
    assert.Equals(t, IsSpider(ua), false)

    ua = "dfs; dfdsf; baiduspider dsf..."
    assert.Equals(t, IsSpider(ua), true)
}

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1312.52 Safari/537.17"

func BenchmarkIsSpider(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsSpider(userAgent)
    }
}
