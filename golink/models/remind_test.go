package models

import (
    "github.com/sdegutis/go.assert"
    "testing"
)

func TestRemind(t *testing.T) {
    var userId int64 = 1
    r, err := Remind_ForUser(userId)
    assert.Equals(t, err, nil)

    err = Remind_Inc(userId, REMIND_COMMENT)
    assert.Equals(t, err, nil)
    err = Remind_Inc(userId, REMIND_COMMENT)
    assert.Equals(t, err, nil)
    err = Remind_Inc(userId, REMIND_FANS)
    assert.Equals(t, err, nil)

    r2, err := Remind_ForUser(userId)
    assert.Equals(t, err, nil)
    assert.Equals(t, r2.Comments, r.Comments+2)
    assert.Equals(t, r2.Fans, r.Fans+1)

    err = Remind_Reset(userId, REMIND_COMMENT)
    assert.Equals(t, err, nil)

    r3, err := Remind_ForUser(userId)
    assert.Equals(t, err, nil)
    assert.Equals(t, r3.Comments, 0)
    assert.NotEquals(t, r3.Fans, 0)
}
