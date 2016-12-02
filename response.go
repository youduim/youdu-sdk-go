package ydapp

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	StatusOK = 0
)

var (
	ErrNoSuchField = errors.New("no such param")
)

type RawMsg struct {
	Data   []byte
	Length int32
	AppId  string
}

type ApiResponse struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	param   map[string]interface{}
	body    []byte
}

func NewReponse(bs []byte) (*ApiResponse, error) {
	rsp := ApiResponse{
		body: bs,
	}

	err := json.Unmarshal(bs, &rsp.param)
	if err != nil {
		return nil, err
	}
	rsp.ErrCode, err = rsp.GetInt32("errcode")
	if err != nil {
		return nil, err
	}
	rsp.ErrMsg, err = rsp.GetString("errmsg")
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (this *ApiResponse) GetString(key string) (string, error) {
	n, ok := this.param[key]
	if !ok {
		return "", ErrNoSuchField
	}
	s, ok := n.(string)
	if !ok {
		return "", errors.New("type assertion to string failed")
	}
	return s, nil
}

func (this *ApiResponse) GetInt32(key string) (int32, error) {
	n, ok := this.param[key]
	if !ok {
		return 0, ErrNoSuchField
	}
	nn, ok := n.(float64)
	if !ok {
		return 0, errors.New("type assertion to float failed")
	}
	return int32(nn), nil
}

func (this *ApiResponse) Status() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", this.ErrCode, this.ErrMsg)
}

func (this *ApiResponse) StatusOK() bool {
	return this.ErrCode == StatusOK
}

func (this *ApiResponse) Error() error {
	if this.StatusOK() {
		return nil
	}
	return errors.New(this.Status())
}

func (this *ApiResponse) Body() []byte {
	return this.body
}
