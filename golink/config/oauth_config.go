package config

import (
    "github.com/QLeelulu/ohlala/golink"
)

type oauthConfig struct {
    ClientId     string
    ClientSecret string
    Scope        string
    AuthURL      string
    TokenURL     string
    RedirectURL  string
}

var OAuthConfigs map[string]*oauthConfig

func init() {
    OAuthConfigs = make(map[string]*oauthConfig)

    oauthRedirectUrl := golink.Host_Name + "/user/oauth2calback?from="

    googleConfig := &oauthConfig{
        ClientId:     "1098296103309.apps.googleusercontent.com",
        ClientSecret: "g707twAeUlECzD4BIy9ShEnD",
        Scope:        "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
        AuthURL:      "https://accounts.google.com/o/oauth2/auth",
        TokenURL:     "https://accounts.google.com/o/oauth2/token",
        RedirectURL:  oauthRedirectUrl + "google",
    }

    OAuthConfigs["google"] = googleConfig
}
