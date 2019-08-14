package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var (
	BaseURL = "https://qyapi.weixin.qq.com"

	Methods = map[string]HTTPAction{
		"GetAccessToken": HTTPAction{
			Method: "GET",
			URL:    "/cgi-bin/gettoken?corpid={{.CropID}}&corpsecret={{.Secret}}",
		},
		"MessageSend": HTTPAction{
			Method: "POST",
			URL:    "/cgi-bin/message/send?access_token={{.Token}}",
		},
	}
)

type HTTPAction struct {
	Method string
	URL    string
}

type API struct {
	baseURL     string
	Secret      string
	AgentID     int
	CropID      string
	Token       string
	tokenExpire time.Time
	clt         *http.Client
}

type RespBase struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type TokenResp struct {
	RespBase
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`
}

func (r *RespBase) Error() error {
	if r.ErrCode == 0 {
		return nil
	}
	return fmt.Errorf("code:%d error:%s", r.ErrCode, r.ErrMsg)
}

type MsgResp struct {
	RespBase
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

func NewAPI(cropID string, agentID int, secret string) *API {
	a := new(API)
	a.CropID = cropID
	a.AgentID = agentID
	a.Secret = secret
	a.baseURL = BaseURL
	a.clt = &http.Client{}
	return a
}

func (a *API) autoGetToken() (err error) {
	if time.Since(a.tokenExpire) > 0 {
		_, err = a.GetToken()
	}
	return
}

func (a *API) GetToken() (token string, err error) {
	var resp TokenResp
	err = a.doHTTP("GetAccessToken", nil, &resp)
	if err != nil {
		return
	}
	err = resp.Error()
	if err != nil {
		return
	}
	a.Token = resp.AccessToken
	a.tokenExpire = time.Now().Add(time.Duration(resp.ExpireIn) * time.Second)
	token = a.Token
	return
}

func (a *API) Send(msg interface{}) (err error) {
	var resp MsgResp
	err = a.doHTTP("MessageSend", msg, &resp)
	if err != nil {
		return
	}
	err = resp.Error()
	if err != nil {
		return
	}
	return
}

func (a *API) SendTextToUser(content string, users ...string) (err error) {
	msg := NewTextMsg(content)
	if len(users) == 0 {
		msg.ToUser = "@all"
	} else {
		msg.ToUser = strings.Join(users, "|")
	}
	msg.AgentID = a.AgentID
	err = a.Send(msg)
	return
}

func (a *API) SendMarkdownToUser(content string, users ...string) (err error) {
	msg := NewMarkdownMsg(content)
	if len(users) == 0 {
		msg.ToUser = "@all"
	} else {
		msg.ToUser = strings.Join(users, "|")
	}
	msg.AgentID = a.AgentID
	err = a.Send(msg)
	return
}

func (a *API) doHTTP(key string, reqBody interface{}, data interface{}) (err error) {
	act, ok := Methods[key]
	if !ok {
		err = errors.New("No such key")
		return
	}
	if key != "GetAccessToken" {
		err = a.autoGetToken()
		if err != nil {
			return
		}
	}
	tmpl := template.New(key)
	tmpl, err = tmpl.Parse(act.URL)
	if err != nil {
		return
	}
	out := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(out, a)
	if err != nil {
		return
	}
	reqURL := a.baseURL + out.String()
	var body *bytes.Buffer
	if reqBody != nil {
		var b []byte
		b, err = json.Marshal(reqBody)
		if err != nil {
			return
		}
		body = bytes.NewBuffer(b)
	} else {
		body = bytes.NewBuffer([]byte{})
	}
	req, err := http.NewRequest(act.Method, reqURL, body)
	if err != nil {
		return
	}
	resp, err := a.clt.Do(req)
	if err != nil {
		return
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(data)
	if err != nil {
		return
	}
	return
}
