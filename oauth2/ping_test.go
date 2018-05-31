package oauth2

import (
	. "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pingTestUserResponse = `{
 "sub": "248289761001",
 "name": "Jane Doe",
 "given_name": "Jane",
 "family_name": "Doe",
 "preferred_username": "jane.doe",
 "email": "janedoe@example.com",
 "picture": http://example.com/janedoe/me.jpg
}`

func Test_Ping_getUserInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Equal(t, "secret", r.FormValue("access_token"))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(pingTestUserResponse))
	}))
	defer server.Close()

	pingAPI = server.URL

	u, rawJSON, err := providerGithub.GetUserInfo(TokenInfo{AccessToken: "secret"})
	NoError(t, err)
	Equal(t, "248289761001", u.Sub)
	Equal(t, "janedoe@example.com", u.Email)
	Equal(t, "Jane Doe", u.Name)
	Equal(t, githubTestUserResponse, rawJSON)
}
