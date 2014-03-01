package crawler

import (
    "errors"
    "log"
    "math/rand"
    "strings"

    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/controllers"
    "github.com/QLeelulu/ohlala/golink/forms"
    "github.com/QLeelulu/ohlala/golink/models"
)

var topTopics []models.Topic
var lenTopTopics int

func init() {
    var err error
    topTopics, err = models.Topic_GetTops(1, 200)
    if err != nil {
        log.Fatalln("load top topics error:", err.Error())
    }
    lenTopTopics = len(topTopics)
}

type Crawler interface {
    Run() error
}

type BaseCrawler struct {
    Name    string
    Url     string
    UserIds []int64
}

func (self *BaseCrawler) saveLink(url, title string) (err error) {
    defer func() {
        if err != nil {
            if strings.Index(err.Error(), "Url已经提交过") > -1 {
                goku.Logger().Logln("Crawler saveLink:", err.Error(), url, title)
            } else {
                goku.Logger().Errorln("Crawler saveLink error:", err.Error(), url, title)
            }
        }
    }()
    idCount := len(self.UserIds)
    if idCount < 1 {
        err = errors.New("no user ids")
        return
    }
    userId := self.UserIds[rand.Int63n(int64(idCount))]
    user := models.User_GetById(userId)
    if user == nil || user.Id < 1 {
        err = errors.New("no selected user")
        return
    }

    if strings.Index(url, "news.dbanotes.net") > 0 {
        return nil
    }
    // 移除多余的字符
    if strings.LastIndex(title, ")") == len(title)-1 && strings.Index(title, " (") > 0 {
        title = title[0:strings.LastIndex(title, " (")]
    }

    topics := []string{}
    ltitle := strings.ToLower(title)
    for i := 0; i < lenTopTopics; i++ {
        if strings.Index(ltitle, topTopics[i].NameLower) > -1 {
            if len(topTopics[i].Name) > 1 {
                topics = append(topics, topTopics[i].Name)
            }
        }
    }

    m := map[string]string{
        "title":   title,
        "context": url,
        "topics":  strings.Join(topics, ","),
    }
    f := forms.CreateLinkSubmitForm()
    f.FillByMap(m)

    success, linkId, errMsg, _ := models.Link_SaveForm(f, user.Id, false)

    if success {
        go controllers.AddLinkForSearch(0, f.CleanValues(), linkId, user.Name) //contextType:0: url, 1:文本   TODO:
    } else {
        err = errors.New(strings.Join(errMsg, ", "))
        return
    }
    return nil
}
