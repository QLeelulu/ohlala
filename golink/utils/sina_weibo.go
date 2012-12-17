package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
//"fmt"
)

type SinaUserInfo struct {
	Screen_Name string `json:screen_name`
}

type SinaWeiBo struct {
	accessToken     AccessToken
}

func NewSinaWeiBo(token AccessToken) *SinaWeiBo {
	return &SinaWeiBo{token}
}


func (s *SinaWeiBo) GetUserInfo() (SinaUserInfo, error) {
	v := url.Values{}
	v.Add("access_token", s.accessToken.Access_Token)
	v.Add("uid", s.accessToken.Uid)
	url := "https://api.weibo.com/2/users/show.json?" + v.Encode()
	response, err := http.Get(url)

	if err != nil {
		return SinaUserInfo{}, err
	}
	defer response.Body.Close()
	jsonMap := SinaUserInfo{}

	json.NewDecoder(response.Body).Decode(&jsonMap)

	return jsonMap, nil
}
