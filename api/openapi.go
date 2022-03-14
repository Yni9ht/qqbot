package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Yni9ht/qqbot/log"
	"github.com/Yni9ht/qqbot/model"
	"io/ioutil"
	"net/http"
)

const (
	domain        = "api.sgroup.qq.com"
	sandboxDomain = "sandbox.api.sgroup.qq.com"
)

type uri string

const (
	WsGateway uri = "/gateway"

	SendMsg uri = "/channels/%s/messages"
)

type OpenAPI struct {
	URL     string
	Token   string
	Sandbox bool
	Client  *http.Client
}

func (o *OpenAPI) GetUrl(endpoint uri, param ...interface{}) string {
	d := domain
	if o.Sandbox {
		d = sandboxDomain
	}
	if len(param) > 0 {
		return fmt.Sprintf("https://%s%s", d, fmt.Sprintf(string(endpoint), param...))
	}
	return fmt.Sprintf("https://%s%s", d, endpoint)
}

func (o *OpenAPI) GetGatewayURL() (string, error) {
	req, err := http.NewRequest(http.MethodGet, o.GetUrl(WsGateway), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", o.Token)
	res, err := o.Client.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if !isSuccessRes(res.StatusCode) {
		return "", errors.New(fmt.Sprintf("code:%v, msg:%v", res.StatusCode, string(b)))
	}

	result := &model.WebsocketGateway{}
	err = json.Unmarshal(b, result)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func (o *OpenAPI) SendMessage(channelID string, msg model.OpenAPIMessageReq) error {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, o.GetUrl(SendMsg, channelID), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", o.Token)
	req.Header.Set("Content-Type", "application/json")
	res, err := o.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if !isSuccessRes(res.StatusCode) {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("code:%v, msg: %v", res.StatusCode, string(body)))
	}

	log.Infof("send message: %v", res)
	return nil
}

var successStatusSet = map[int]bool{
	http.StatusOK:        true,
	http.StatusNoContent: true,
}

func isSuccessRes(code int) bool {
	if _, ok := successStatusSet[code]; ok {
		return true
	}
	return false
}
