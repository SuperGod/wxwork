package api

import (
	"testing"
)

func getAPI(t *testing.T) *API {
	a, err := GetAPI()
	if err != nil {
		t.Fatal(err.Error())
	}
	return a
}

func TestDoHttp(t *testing.T) {
	api := getAPI(t)
	var data TokenResp
	err := api.doHTTP("GetAccessToken", nil, &data)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data)

}

func TestGetToken(t *testing.T) {
	api := getAPI(t)
	token, err := api.GetToken()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(token)
	t.Log(api.tokenExpire)
}

func TestSendText(t *testing.T) {
	api := getAPI(t)
	msg := NewTextMsg("hello, I am a bot")
	// msg.ToParty = "1"
	msg.ToUser = "@all"
	msg.AgentID = api.AgentID
	err := api.Send(msg)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSendTextToUser(t *testing.T) {
	api := getAPI(t)
	err := api.SendTextToUser("I am a bot msg")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSendMarkdownToUser(t *testing.T) {
	api := getAPI(t)
	err := api.SendMarkdownToUser(`# bot
## I am a bot msg`)
	if err != nil {
		t.Fatal(err.Error())
	}
}
