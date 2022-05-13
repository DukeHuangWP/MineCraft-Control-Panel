package googleOauth2

import (
	"minecraft-control-panel/internal/global"

	"golang.org/x/oauth2"
	googleInc "golang.org/x/oauth2/google"
)

var (
	OauthConfig *oauth2.Config
)

//Google帳號資訊
type GoogleAcc struct {
	ID            int
	Email         string
	VerifiedEmail bool
	PictureUrl    string
}

func init() {
	//https://console.cloud.google.com/apis/credentials
	var redirectURL string
	if global.HostDomainName == "localhost" || global.HostDomainName == "127.0.0.1" || global.HostDomainName == "0.0.0.0" {
		redirectURL = global.HostScheme + global.HostDomainName + ":" + global.HostPort + "/" + global.WebURLRoot + "/" + global.Oauth2CallbackName

	} else {
		redirectURL = global.HostScheme + global.HostDomainName + "/" + global.WebURLRoot + "/" + global.Oauth2CallbackName
	}
	OauthConfig = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     global.Oauth2ClientID,
		ClientSecret: global.Oauth2ClientSecret, //set from google credentials
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     googleInc.Endpoint,
	}
}
