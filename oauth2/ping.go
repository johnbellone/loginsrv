package oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/loginsrv/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type PingUser struct {
	UserID            string `json:"sub,omitempty"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	Picture           string `json:"picture,omitempty"`
}

var pingProvider Provider

func init() {
	pingProvider = Provider{
		Name:     "ping",
		AuthURL:  "https://%v/as/authorization.oauth2",
		TokenURL: "https://%v/as/token.oauth2",
		GetUserInfo: func(token TokenInfo) (model.UserInfo, string, error) {
			pu := PingUser{}

			// Use the provider's domain to for the IdP domain; this is set from configuration.
			u, err := url.Parse(pingProvider.AuthURL)
			if err != nil {
				return model.UserInfo{}, "", err
			}

			ue := fmt.Sprintf("https://%v/idp/userinfo.openid?schema=openid&access_token=%v", u.Host, token.AccessToken)
			resp, err := http.Get(ue)
			if err != nil {
				return model.UserInfo{}, "", err
			}

			if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
				return model.UserInfo{}, "", fmt.Errorf("wrong content-type on ping get user info: %v", resp.Header.Get("Content-Type"))
			}

			if resp.StatusCode != 200 {
				return model.UserInfo{}, "", fmt.Errorf("got http status %v on ping get user info", resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return model.UserInfo{}, "", fmt.Errorf("error reading ping get user info: %v", err)
			}

			err = json.Unmarshal(b, &pu)
			if err != nil {
				return model.UserInfo{}, "", fmt.Errorf("error parsing ping get user info: %v", err)
			}

			return model.UserInfo{
				Sub:     pu.UserID,
				Picture: pu.Picture,
				Name:    pu.Name,
				Email:   pu.Email,
				Origin:  "ping",
			}, string(b), nil
		},
	}

	RegisterProvider(pingProvider)
}
