package test

import (
	"chat/cmd/chat/api"
	"chat/cmd/chat/util"
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenLoginParam(t *testing.T) {
	clientApiVersion := "0.0.1"

	param := api.LoginParam{
		RequestHeader: api.RequestHeader{
			Version:   clientApiVersion,
			MsgID:     util.GenerateMsgID(),
			Method:    api.Login,
			Timestamp: util.GenerateTimestamp(),
		},
		AccountName: "LeeBai",
		AccessToken: "LeeBai-Token",
	}

	request, _ := json.Marshal(param)
	fmt.Printf("%v\n", string(request))
}

func TestGenLogoutParam(t *testing.T) {
	clientApiVersion := "0.0.1"

	param := api.LogoutParam{
		RequestHeader: api.RequestHeader{
			Version:   clientApiVersion,
			MsgID:     util.GenerateMsgID(),
			Method:    api.Logout,
			Timestamp: util.GenerateTimestamp(),
		},
		AccountName: "LeeBai",
		AccessToken: "LeeBai-Token",
	}

	request, _ := json.Marshal(param)
	fmt.Printf("%v\n", string(request))
}

func TestGenCreateGroupParam(t *testing.T) {
	clientApiVersion := "0.0.1"

	param := api.CreateGroupParam{
		RequestHeader: api.RequestHeader{
			Version:   clientApiVersion,
			MsgID:     util.GenerateMsgID(),
			Method:    api.CreateGroup,
			Timestamp: util.GenerateTimestamp(),
		},
		SenderID: "LeeBai",
		GroupID:  "group1",
	}

	request, _ := json.Marshal(param)
	fmt.Printf("%v\n", string(request))
}

func TestGenDestroyGroupParam(t *testing.T) {
	clientApiVersion := "0.0.1"

	param := api.DeleteGroupParam{
		RequestHeader: api.RequestHeader{
			Version:   clientApiVersion,
			MsgID:     util.GenerateMsgID(),
			Method:    api.DeleteGroup,
			Timestamp: util.GenerateTimestamp(),
		},
		SenderID: "LeeBai",
		GroupID:  "group1",
	}

	request, _ := json.Marshal(param)
	fmt.Printf("%v\n", string(request))
}

func TestGenKeepAliveParam(t *testing.T) {
	clientApiVersion := "0.0.1"

	param := api.AppKeepAliveParam{
		RequestHeader: api.RequestHeader{
			Version:   clientApiVersion,
			MsgID:     util.GenerateMsgID(),
			Method:    api.Keepalive,
			Timestamp: util.GenerateTimestamp(),
		},
	}

	request, _ := json.Marshal(param)
	fmt.Printf("%v\n", string(request))
}
