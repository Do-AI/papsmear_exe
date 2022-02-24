package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// user 서비스에게 agent가 유효한지를 체크해주는 로직이 들어있다.

var authInstance *auth
var once sync.Once

type login struct {
	TokenType    string      `json:"token_type"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	Scope        interface{} `json:"scope"`
	Expire       int         `json:"expires_in"`
	Jti          string      `json:"jti"`
}

type auth struct {
	tokenExpireAt int
	headerValue   string
}

func getAuthInstance() *auth {
	once.Do(func() {
		authInstance = &auth{0, ""}
	})
	return authInstance
}

func (a *auth) signIn() *login {
	endpoint := CONFIG.Svc.APIURL + "/user/oauth/token"

	client := &http.Client{}

	bodyData := url.Values{
		"client_id":     {CONFIG.Credential.ClientID},
		"client_secret": {CONFIG.Credential.ClientSecret},
		"username":      {CONFIG.Credential.Username},
		"grant_type":    {"password"},
	}

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBufferString(bodyData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var respJson login
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Fatal(err)
	}

	return &respJson
}

func GetAuthHeaderValue() string {
	instance := getAuthInstance()
	if instance.tokenExpireAt > int(time.Now().Unix()) {
		return instance.headerValue
	}

	AuthValue := instance.signIn()
	instance.headerValue = fmt.Sprintf("%s %s", AuthValue.TokenType, AuthValue.AccessToken)
	instance.tokenExpireAt = (AuthValue.Expire-30)*1000 + int(time.Now().Unix())
	return instance.headerValue
}
